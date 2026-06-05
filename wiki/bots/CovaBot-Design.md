# CovaBot (Go Rewrite) — Design Record

> Living document. Status tags: **[DECIDED]** = settled, **[OPEN]** = needs a call, **[IDEA]** = parked, not committed.
> Last updated: 2026-06-05

---

## 0. Goal

Rewrite CovaBot in Go as a Discord personality bot that reads as a *live user*: remembers people, holds opinions that evolve, responds with the cadence and selectivity of a real person — not an eager assistant that replies to everything or a rate-limited one that ghosts mid-conversation.

---

## 1. Core Organizing Principle  **[DECIDED]**

Separate **what Cova *is*** from **when and how Cova *speaks***.

- **Prompt files** own identity, voice, preferences, conversational nuance. Declarative, version-controlled, editable without recompiling.
- **Go code** owns orchestration: whether to engage the LLM, which prompt fragments to assemble, how much context to attach, reliability of output.

Avoid leakage in either direction:
- Personality logic in Go → tuning Cova means a redeploy. Bad.
- Orchestration logic in prompts → the LLM is asked to both *be* Cova and *decide whether to respond* in one muddled blob. This was the buried-IGNORE-marker problem in the TS version.

---

## 2. Model Tiers  **[DECIDED]**

Code requests a **capability tier**, not a specific model. Provider config maps tier → actual model, layered on the existing provider priority (Ollama → Anthropic → Gemini → OpenAI). Net result is a **tier × provider matrix**.

- **Low** — relevance gate; LLM fallback for conversation tagging. High-frequency, latency-sensitive, deliberately lean context (no full voice demo).
- **Med** — background reasoning, not user-facing, runs at human pace (conversations die slowly, so the bigger model is affordable). Jobs:
  - **Stance evolution + formation** — on salient conversation death, reconcile Cova's stances: *update* existing ones and *create new* ones (he argues about something he had no prior opinion on → new stance row).
  - **Conversation summarization** — compress a dead salient conversation into a retrievable episodic-memory record.
  - **Memory maintenance (forgetting)** — periodic sweep: merge near-duplicate episodic records, decay/drop old low-salience ones. Long-term analogue of the live-set bound (§6.4); keeps retrieval sharp as the store fills.
  - Triggers: event-driven on conversation death + a periodic maintenance sweep.
  - *Not med:* IDF tag-frequency counts (§6.2) are a cheap DB rollup, no LLM call.
- **High** — response generation. Full voice + full context budget. Runs least often.

---

## 3. Decision Pipeline (code side)  **[DECIDED]**

Staged, cheap-to-expensive. Each stage returns `respond` / `abstain` / `defer`.

```
message
  → hard filter            (own msgs, ignored channels — microseconds)
  → conversation tagging   (mechanical-first, low-LLM fallback)
  → pull score             (membership + mention + stance; restraint modulates)
  → [pull below floor?] → abstain
  → relevance gate         (LOW tier, lean context) → yes/no + which convo(s) + reason
  → [no?] → abstain
  → generation             (HIGH tier, full context; seeded with gate's reason)
```

Ordering rule: never spend an expensive step on something a cheap step already settled.

Two-call rationale: separating gate from generation lets each call get *different* context (lean for the gate, full for generation), is more reliable than a folded IGNORE marker, and is nearly free given local Ollama for the low tier.

### 3.1 Gate → generation contract  **[DECIDED]**

**Gate's role:** *not* a from-scratch relevance judge — pull already did the heuristic work and the pull-floor decides whether the gate even runs. The gate (a) catches false positives pull can't (rhetorical questions, sarcasm, name quoted unrelatedly, redundant responses) and (b) articulates intent to seed generation. Lean context — no full voice demo.

```
GateResult {
  respond:         bool
  conversation_id: which thread the response belongs to   // code proposes top-pull; gate may override
  reason:          terse string — WHY he's responding      // seeds generation
  energy:          {quick-jab | normal | invested}         // drives cadence/length
}
```

`reason` + `energy` flow into the generation prompt as **natural framing, not a status block** — the fix for the old TS problem where engagement context read like a checklist. Generation gets *"Someone's pushing back on your KH canon claim — you've got a strong take, fire back, keep it short,"* not `mentioned: NO / last-responded: 120s`. `energy` is what produces human cadence (a low-pull stance butt-in → one-liner; an invested thread → full reply).

Binary respond/abstain for v1. A "defer / let someone else take it" option is a possible later cadence refinement, not now.

---

## 4. Responsiveness Model  **[DECIDED]**

**Root cause of most symptoms:** no conversational state, plus a "social battery" that only models *restraint* and is allowed to veto everything — including direct questions.

Two axes:

- **Pull** (per message): how much this message warrants a response.
  - Direct mention → max
  - Reply to something Cova just said → very high
  - Topic Cova has a stance on → elevated
  - Random chatter → low
  - Computed as max over the conversations the message is tagged to, weighted by Cova's engagement in each.
- **Restraint** (the old battery): dominating? channel pace? recency of his last message.

**The rule that makes him feel alive:** restraint *modulates* low-pull messages but **never vetoes high-pull ones.**

Fixes: responds-to-everything, doesn't-respond-when-he-should, and the question→answer→followup→silence drop (followup is highest-pull, was being killed by depleted battery).

### 4.1 Pull combination — MAX, not blend  **[DECIDED]**

> **Note for the implementing agent:** this was a deliberate choice between two structures. Read the tradeoff before changing it.

The three signal families combine by **max**, not by weighted sum:

```
pull = max(intrinsic, conversational, stance)
```

i.e. the single strongest reason to respond wins. A weak signal never accumulates with other weak signals to clear the bar.

**Why max (chosen):**
- Matches how a real person responds — you reply when there's *one clear reason*, not when three faint nudges quietly add up.
- Trivially debuggable: every response traces to one dominant signal ("responded because it was a direct mention").
- Structurally resists the **responds-to-everything** failure mode, which is the #1 behavior we are trying to eliminate.

**Blend (weighted sum) — considered and REJECTED:**
- Would produce more *ambient* responsiveness (a weakish reply, in a semi-warm conversation, on a mild-interest topic could sum past the threshold even though no single signal would).
- Rejected because it is hard to debug ("why did he respond to *that*?") and drifts back toward chatty-about-everything unless the weights are kept very tight — reintroducing the exact bug we're fixing.

**Guardrail for the agent:** do **not** switch to a blend/weighted-sum to make Cova more talkative. If more ambient chime-in behavior is wanted, the correct lever is to **raise the `stance` term** (so "topic Cova cares about" alone clears the bar) — not to let signals accumulate. Keep the combination as max.

### 4.2 Pull values — affinity + stance interaction  **[DECIDED]**

**Participation affinity** scales conversational pull by Cova's status in a conversation. Starting values (tunable; the *ordering* is the fixed part):

| Status   | Meaning                              | Affinity |
|----------|--------------------------------------|----------|
| in       | actively in it, recent               | ~1.0     |
| was-in    | spoke earlier, gone cold             | ~0.5     |
| never    | never spoke in it                    | ~0.15    |

`was-in` is set so a reply to his old point pulls him back, but he isn't glued to dead threads.

**Stance × membership — SCALED.** Because pull is a `max` (§4.1), stance pull already bypasses membership for free — there is **no separate "stance overrides membership" rule**, and the agent should not add one. The only rule: stance pull is computed everywhere, but **scaled down in `never` conversations** so that only Cova's *strongest* stances clear the bar there.

- Effect: he'll butt into a stranger's thread, but only on a hill he'd die on — not a mild preference.
- This is what delivers the organic "spots a Superman thread across the channel and dives in" behavior **without** turning him into a channel-wide "well, actually" pest.

### 4.3 Global activity level — the volume knob  **[DECIDED]**

An operator-facing master volume, built as a modifier on the **pull floor** (not a parallel system):
```
effective_floor = base_floor + channel_baseline_offset + temporary_dampener(decaying)
```

- **Baseline** — persistent, operator-set resting floor. **Per-channel with a global default** (chatty in memes, reserved in serious channels).
- **Temporary dampener** — the "quiet him for a while" knob. Adds a large floor offset that **decays back to baseline** over a set duration — pulls back, then self-recovers, no need to remember to un-quiet.
- **Mute** — hard floor; only owner/admin direct address passes (so you can still reach him to flip it back).

**Critical property (the lesson from the old social battery):** raising the floor suppresses **ambient chime-in only** — direct address (mention / reply / `addressee==cova`) still clears it, *except* in full mute. "Quieted" reads as a reserved real person, not a bot ignoring you when you talk to it.

**Triggers — both:**
- *Manual* — command sets baseline or fires the temporary dampener.
- *Automatic* — negative feedback (people saying "stop"/"shut up covabot", negative reactions on his messages) fires the temporary dampener. "Reads the room and dials it back" — the live-user behavior the project is chasing.

---

## 5. Engagement State  **[DECIDED]**

First-class object, per channel. The missing primitive.

- Tracks who Cova is talking to, how recently, the active thread.
- Active engagement drops the pull threshold near zero for thread participants, then decays as the thread cools.
- Produces human cadence (leans in, drifts off) and feeds the generation prompt so Cova knows he's mid-conversation → fixes "says hello like he's never seen anyone."

---

## 6. Conversations — Many-to-Many Membership  **[DECIDED]**

Channels are interleaved. Rather than partitioning messages into one conversation each, **a message is tagged with the set of conversations it applies to**, weighted by confidence. Overlap is a feature, not a problem — this sidesteps hard disentanglement entirely.

Data shape:
```
message_conversation(message_id, conversation_id, weight)   -- join, many-to-many
```

- A **conversation** = the set of messages tagged to it. Participants, recency, and topic-centroid are **derived** from its members — nothing extra to keep in sync. No recent members → it ages out (lossy, decaying).
- **Cova-membership** (in / was-in / never) is derived: did Cova author any tagged message, how recently.
- **Bridging:** a message tagged to two conversations is how a topic migrates between threads — and the hook for Cova noticing his thread brushed a topic he has opinions about.

**Tagging mechanic — mechanical-first, LLM-fallback:**
1. Reply references (free, definitive)
2. Embedding similarity to each live conversation's centroid (cheap via pgvector)
3. Participant overlap + recency (free)
4. Escalate to LOW-tier LLM *only* for ambiguous "new thread vs topic shift" calls.

Uncertain membership → default to channel-general (degrade to baseline, never confidently misfire).

### 6.1 Tag scheme — two namespaces  **[DECIDED]**

Every message gets tagged. Tags fall into two namespaces with different rules, because they drive different things:

- **Topical tags** — open vocabulary, LLM-generated, **embedded and fuzzy-matched** (never string-matched). Multi-granularity: broad → entity → specific (e.g. `kingdom hearts` / `donald duck` / `zettaflare`). These drive *conversation membership* and *stance matching*. They never touch code logic, so they don't need to be canonical — the LLM can phrase the same idea differently across messages and embedding similarity still clusters them.
- **Structural tags** — closed vocabulary, defined by us, **exact-matched enums**. Code branches on these, so they must be reliable.

  *Filter for inclusion:* a structural tag earns a slot only if (a) code actually branches on it **and** (b) it can't be derived more cheaply (from Discord metadata or from the clustering step itself).

  **From Discord metadata (code, no LLM):**
  - `direct-mention` — Cova was @mentioned
  - `reply-to-cova` — Discord reply pointing at a Cova message. *High-precision, low-recall — most people don't use the reply feature, so this is a strong bonus when present but NOT load-bearing. See directedness note below.*

  **From the tagger (LLM, closed enum, judged with thread context):**
  - `addressee` ∈ `{cova, other-user, room}` — `cova` = semantically addressed to him even without a mention; `other-user` = clearly someone else's exchange (suppress); `room` = ambient
  - `intent` ∈ `{question, statement, low-effort}` — `low-effort` (greetings, "lol", reactions) damps pull

  **Explicitly cut** (derived elsewhere, not declared by the tagger): `cova-stance-hit` (code computes from topical tags × stance store, §6.2), `conversation-opener` (a message matching no live conversation *is* an opener), `sentiment` (no hot-path branch; med-tier derives it during stance evolution if needed).

Rule of thumb: anything an LLM free-generates is topical; anything code makes a decision on is structural.

Tagging is done by a persona-free classification pass (separate from generation — the tagger reads the conversation neutrally and emits tags, it does not respond as Cova). Don't ask the LLM for structural facts Discord metadata already provides (`direct-mention`, `reply-to-cova` are free from message metadata).

**Tag Cova's own responses too** — that's how his participation registers (`in`/`was-in`) and how his stances accumulate.

**Dedup before use:** the tagger over-fragments (`zettaflare`, `100x kamehameha`, `beam struggle`, `kamehameha zettaflare beam struggle` ≈ one concept). Collapse near-identical topical tags in embedding space before weighting, or they get multi-counted.

**Directedness is mostly a conversation-state property, not a per-message one.** Because the reply feature is rarely used, the system can't lean on `reply-to-cova`. The common-case workhorse is **engagement-continuity** (§5): when Cova just spoke, he's `in` the hot thread, so a replyless, address-less follow-up from a participant still pulls a response. Directedness stacks as:
```
directedness = max(reply-to-cova [rare, near-certain],
                   addressee==cova [inferred w/ thread context],
                   engagement-continuity [common-case workhorse])
```
This is scoped to threads Cova is already in, so it does not recreate responds-to-everything; restraint still modulates if he's been dominating.

### 6.2 Tag specificity — broad vs specific  **[DECIDED]**

Broad and specific tags are load-bearing for *different* jobs and pull in opposite directions. (Observed in a real thread: broad tags like `kingdom hearts` / `kh vs dbz` persist across the whole conversation; specific tags churn message-to-message as the micro-topic drifts.)

- **Broad tags = continuity.** They keep a drifting conversation clustered as one thread (Zettaflare → Riku's comic canon is still one argument). **Membership/clustering must lean on broad/common tags** or a topic-drift looks like a brand-new conversation and Cova loses the thread.
- **Specific tags = precision.** They tell you whether *this exact message* hits something Cova cares about. **Stance/relevance pull must lean on specific/rare tags** (`canon`+`manga` is precise enough to fire a stance; `kingdom hearts` is too coarse).

**Mechanism — IDF, not LLM self-rating.** Derive specificity from observed tag frequency across the corpus (rare = high weight, common = downweighted). Never have the LLM rate how "specific" a tag is.

- Self-calibrates per server: `goku` is common/uninformative in a DBZ channel, rare/informative in a KH channel — the data decides, no hardcoding.
- Apply asymmetrically: **membership** keeps common tags' weight (continuity); **stance/relevance** IDF-weights hard (a `zettaflare` match counts far more than a `kingdom hearts` match).

Net pull formula substrate:
```
conversational_pull = participation_affinity(in/was-in/never) × tag_similarity_to_conversation   // continuity-weighted
stance_pull         = max over Cova's stances of  idf_weighted_similarity(message_tags, stance_subject)  // specificity-weighted
pull                = max(intrinsic, conversational_pull, stance_pull)   // §4.1
```

### 6.3 Membership assignment & thresholds  **[DECIDED]**

What gets embedded for membership: **topical tags**, not raw message text. Broad tags are explicit glue that keeps a drifting thread clustered when the wording changes completely ("yeah but that's not canon" and "imagine Zettaflare blocking a Kamehameha" share no words but share `kingdom hearts`/`kh vs dbz`). Raw-text similarity was **rejected** because it splits one argument into two conversations exactly on the lexically-divergent continuations this design exists to handle. Consequence: the LLM tagger runs on basically every substantive message — acceptable on local Ollama; raw-text is the fallback only if local throughput becomes the bottleneck.

Assignment for an incoming message (post hard-filter):

1. **Reply reference → automatic membership.** Discord reply pointing into conversation X = member of X, bypass the similarity math. Rare but definitive.
2. **Else embed topical tags, compare to each live conversation centroid:**
   - `sim ≥ τ_high` → member, `weight = sim`. Mechanical, no LLM.
   - `τ_low ≤ sim < τ_high` → ambiguous candidate.
   - `sim < τ_low` for *every* live conversation → belongs to nothing → **seed a new conversation**, centroid = this message.
3. **Ambiguous band → ONE batched low-tier call** (not per-candidate): "message + these K candidate conversations — which does it belong to?" Returns the membership set. LLM runs only on genuinely uncertain cases.

- **Thresholds:** start ~`0.75` / ~`0.45` cosine. *Tunable* — right values depend on the embedding model; calibrate against real channel data. The fixed part is the three-zone shape.
- **Centroid** = running mean of member messages' tag-embeddings. Young conversations (one message) have noisy centroids → expect more middle-band escalation early in a thread's life.
- **Bridging** falls out: a message above `τ_high` for two centroids joins both, no special case.

### 6.4 Decay & consolidation  **[DECIDED]**

Framing: **conversations are working memory** (ephemeral, bounded, decaying); **episodic + stance stores are long-term memory** (durable). Death is not deletion — it's the *consolidation* step that moves salient content from the former into the latter. This is a med-tier job.

- **Death trigger — recency TTL.** Each conversation tracks last-activity (newest member message). After a quiet window (start ~30 min, tunable) it's dead and leaves the live set. Adaptive TTL (scaled to the conversation's own cadence) is a possible later refinement, not v1.
- **Live-set bound — TTL + hard top-N cap per channel** (backstop against a pathologically busy channel). Per-channel scoped, so naturally bounded by channel activity. Keeps membership matching (§6.3) fast and free of stale spurious matches.
- **Death ≠ memory write.** Junk conversations (`lol`/`nice`/`yeah`) just evaporate — drop from the live set, no consolidation. Only *salient* conversations get summarized.
- **Salience bar — both modes, toggle, default participated:**
  - *Default (participated):* conversations Cova was in (`in`/`was-in`) → med-tier summarizes into episodic memory + reconciles stances. He loses little vs. witnessed, since stance-pull means he participates in most things he has opinions about anyway.
  - *Toggle (witnessed):* **also** consolidate substantive, stance-adjacent conversations he did *not* participate in — he forms impressions from conversations happening around him. More ambient/human, more med-tier load. Off by default; flip on once the core is stable.
  - Everything else evaporates.
- **Revival needs no special-casing.** A topic resurfacing after its conversation aged out just seeds a fresh conversation (§6.3) — it won't match a dead one. Continuity comes from durable memory retrieval, not from keeping conversation objects immortal. This is the payoff of the working/long-term split: the ephemeral object can die freely because the durable store carries what matters.

---

## 7. Memory — Two Stores  **[DECIDED]**

- **Episodic recall** — pgvector over past exchanges. "What have I said to/about this person before." Retrieved by similarity.
- **Stances** — explicit, *mutable* rows keyed by subject (person/topic), with sentiment + note. Hand-editable. Opinions that evolve live here, NOT as append-only vector soup. The "or at least can be altered" knob. Updated by the MED tier.

---

## 8. Tone Problem  **[DECIDED — approach, not yet implemented]**

Separate from the state problems. Default LLM pleasantness is RLHF-baked and sticky; adjectives won't beat it.

- Demonstrate voice in **example exchanges**, incl. Cova being neutral, dismissive, or bored — not just labeled "opinionated."
- Explicitly grant permission to disagree, to not engage warmly, to give a flat reaction.
- Model choice matters — A/B local models for baseline sycophancy.

---

## 9. Prompt Layer  **[DECIDED]**

Named, single-purpose fragments composed at runtime by a Go assembler from a `PromptSpec` (which fragments, what order, runtime injects). `go:embed`, with optional disk hot-reload for dev.

### 9.1 Fragment taxonomy — four consumers

There is no single prompt; each LLM call needs different material:

| Consumer | Tier | Needs |
|----------|------|-------|
| Tagger | low | tagging instructions + structural enum defs + granularity/dedup guidance. **Persona-free.** |
| Gate | low | relevance-judgment instructions, false-positive cases, `GateResult` contract. Persona-*light* (one-line identity brief max). |
| Generation | high | full voice + register + injects. Where the personality lives. |
| Med jobs | med | task instructions (`stance_reconcile`, `summarize`) + enough identity to judge opinion shifts. |

Fragments split two ways:

- **Persona fragments** (*what Cova is*, reused at different fidelity): `identity` (brief + full variants), `voice`, `preferences`.
- **Task fragments** (per-consumer instructions): `tagging`, `gate`, `generation`, `stance_reconcile`, `summarize`.

The `PromptSpec` per call declares fragments + order + runtime injects. Gate spec = `identity(brief) + gate`; generation spec = `identity(full) + voice + preferences + generation + injects`. Multi-fidelity reuse is what keeps the gate lean and generation rich.

### 9.2 Voice — hybrid

`voice` = **static spine + dynamic retrieval.** A fixed hand-written core (speech patterns + example exchanges incl. neutral/dismissive/bored per §8) anchors consistency; pgvector retrieves Cova's *actual past messages* relevant to the current topic as few-shot examples. Static alone feels canned; pure dynamic drifts and resurfaces weird old messages — hybrid gives consistency *and* "that's actually how he'd say it."

### 9.3 Per-person — three layers, in precedence order

The persona is **one constant voice**; what changes per person is layered on top. Precedence:

1. **Theatrical bit** (§9.4) — explicit, authored, userID-keyed performance overlay. Active when triggered; rides on top of the voice, never replaces it.
2. **Relational register** — automatic, emergent. Inject the addressee's stance row (sentiment/familiarity → warmth, guardedness, teasing) + relevant episodic history (in-jokes, callbacks) into generation. Stance intensity self-scales how much register shifts. Reuses §7; no new subsystem.
3. **Core voice** (§9.2) — the constant default, always underneath.

**Guardrail:** the core voice is **never replaced**. Automatic layers modulate relational framing; bits add a self-aware theatrical overlay — but Cova's baseline voice and wit are always present underneath. Nothing swaps the voice out. (A person's stance row carries a bit more than a topic's: a familiarity/relationship field alongside sentiment + notes.)

### 9.4 Theatrical bits — Cova putting on a voice

**Still Covabot — not a separate persona.** These are *bits*: Cova knowingly hamming it up, cartoony and theatrical, with his own wit and self-awareness showing through. The comedy is "Cova doing an exaggerated voice," NOT "Cova replaced by another character." Always an **overlay on the core voice (§9.2), never a swap** — his baseline stays underneath.

Each bit's fragment is a **performance directive**, not a character sheet — "lay the reverence on absurdly thick, go full Renaissance-faire," not "you are a knight." The exaggeration is the content. Built as a small authored map (a few fragments + userID triggers), not a general per-user-persona feature.

```
Bit {
  trigger_user_id
  trigger_condition:  {present | addressed-or-about}
  overlay_fragment:   bit_*.md          // theatrical overlay on core voice — NEVER a swap
  pull_modifier:      optional
  energy_default:     optional
}
```

Locked examples:
- **Knight** (lord's userID): `trigger: present`; `overlay: bit_knight.md` (high-English, cartoonishly thick reverence — Cova performing a fawning vassal); `pull_modifier: negative sentiment toward the lord → high pull` (leaps to defend even benign slights — a behavioral trigger, not just a voice).
- **Rival** (rival's userID): `trigger: addressed-or-about`; `overlay: bit_contempt.md` (terse, taciturn, occasionally sarcastic — Cova performing exaggerated can't-stand-this-guy); `energy_default: quick-jab`. Intensity ("rival, not mortal enemy") is bounded in the fragment.

**Implementation note — slight detection:** the knight's defend-modifier needs negative-sentiment-toward-the-lord, a per-message judgment. Keep it off the global hot path by computing sentiment **only when a defend-type bit is live in the conversation** (lord present). Bounded to the bit.

---

## 10. Context Budget  **[IDEA — sketched]**

Token budget filled by priority; trim from the bottom when tight.

1. Identity + voice (near-fixed, always)
2. Current trigger message + natural engagement framing (gate's `reason`/`energy`, §3.1)
3. Active persona override / relational register injects (§9.3–9.4)
4. Rolling channel window (last N messages — **with usernames**, the dehumanizing gap from TS)
5. pgvector-retrieved memories + relevant stances + dynamic voice examples

Retrieval is lowest to *include*, highest leverage for *not sounding generic*. Trim history before identity.

---

## 11. Build Sequence  **[DECIDED]**

1. Engagement state + pull/restraint model — fixes the daily responsiveness annoyance.
2. Conversation tagging layered on top to sharpen pull.
3. Memory (episodic + stance) — makes him feel like he knows people.

Don't make conversation tagging a prerequisite for the responsiveness fix.

---

## 12. Open Questions  **[OPEN]**

- **Context budget numbers:** actual per-section token allocations (§10). Pure tuning against the chosen models' context windows — calibrate during build, not a design fork.
