package engagement_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andrewgari/starbunk-go/internal/covabot/engagement"
)

var _ = Describe("Engagement Manager", func() {
	var manager *engagement.Manager
	var chanID = "chan1"
	var userID = "user1"

	BeforeEach(func() {
		manager = engagement.NewManager()
	})

	Describe("Pull Math (max(intrinsic, conversational, stance))", func() {

		Context("Direct Mentions", func() {
			It("should respond to a direct mention with high energy", func() {
				res := manager.ShouldRespond(engagement.MessageInput{
					ChannelID:   chanID,
					AuthorID:    userID,
					IsMentioned: true,
				})
				Expect(res.Respond).To(BeTrue())
				Expect(res.Reason).To(Equal(engagement.ReasonMention))
				// Either normal or invested, but high pull usually means invested or normal
				Expect(res.Energy).To(Or(Equal(engagement.EnergyInvested), Equal(engagement.EnergyNormal)))
			})

			It("should respond to direct mentions even if dampened", func() {
				manager.Dampen(chanID)
				res := manager.ShouldRespond(engagement.MessageInput{
					ChannelID:   chanID,
					AuthorID:    userID,
					IsMentioned: true,
				})
				Expect(res.Respond).To(BeTrue())
				Expect(res.Reason).To(Equal(engagement.ReasonMention))
			})

			It("should respond to direct mentions even if muted", func() {
				manager.SetMute(chanID, true)
				res := manager.ShouldRespond(engagement.MessageInput{
					ChannelID:   chanID,
					AuthorID:    userID,
					IsMentioned: true,
				})
				Expect(res.Respond).To(BeTrue())
				Expect(res.Reason).To(Equal(engagement.ReasonMention))
			})
		})

		Context("Reply to Cova", func() {
			It("should respond to direct replies", func() {
				res := manager.ShouldRespond(engagement.MessageInput{
					ChannelID:   chanID,
					AuthorID:    userID,
					IsReplyToMe: true,
				})
				Expect(res.Respond).To(BeTrue())
				Expect(res.Reason).To(Equal(engagement.ReasonReply))
			})

			It("should respond to replies even if dampened", func() {
				manager.Dampen(chanID)
				res := manager.ShouldRespond(engagement.MessageInput{
					ChannelID:   chanID,
					AuthorID:    userID,
					IsReplyToMe: true,
				})
				Expect(res.Respond).To(BeTrue())
			})

			It("should NOT respond to replies if fully muted", func() {
				manager.SetMute(chanID, true)
				res := manager.ShouldRespond(engagement.MessageInput{
					ChannelID:   chanID,
					AuthorID:    userID,
					IsReplyToMe: true,
				})
				// Full mute only allows direct mentions from admin, simplify to direct mention
				Expect(res.Respond).To(BeFalse())
			})
		})

		Context("Engagement Continuity (Conversational Pull)", func() {
			It("should follow up if Cova recently spoke in the thread", func() {
				manager.RecordCovaSpeak(chanID)
				res := manager.ShouldRespond(engagement.MessageInput{
					ChannelID: chanID,
					AuthorID:  userID,
				})
				Expect(res.Respond).To(BeTrue())
				Expect(res.Reason).To(Equal(engagement.ReasonContext))
			})

			It("should abstain if Cova has not spoken recently", func() {
				res := manager.ShouldRespond(engagement.MessageInput{
					ChannelID: chanID,
					AuthorID:  userID,
				})
				Expect(res.Respond).To(BeFalse())
			})
		})
	})

	Describe("Restraint Model (Dampening)", func() {
		It("dampener should suppress ambient chime-in and engagement continuity", func() {
			manager.RecordCovaSpeak(chanID)
			manager.Dampen(chanID)

			res := manager.ShouldRespond(engagement.MessageInput{
				ChannelID: chanID,
				AuthorID:  userID,
			})
			Expect(res.Respond).To(BeFalse())
		})
	})
})
