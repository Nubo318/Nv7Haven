package discord

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) specials(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "rob") {
		b.checkuser(m)

		if !(len(m.Mentions) > 0) {
			s.ChannelMessageSend(m.ChannelID, "You need to mention the person you are going to rob!")
			return
		}

		exists, suc := b.exists(m, "currency", "user=?", m.Mentions[0].ID)
		if !suc {
			return
		}
		if !exists {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User <@%s> has never used this bot's currency commands.", m.Mentions[0].ID))
			return
		}

		user1, suc := b.getuser(m, m.Author.ID)
		if !suc {
			return
		}

		ups, exists := user1.Properties["rob"]
		if !exists {
			s.ChannelMessageSend(m.ChannelID, "You need property `rob` to Rob people!")
			return
		}

		user2, suc := b.getuser(m, m.Mentions[0].ID)
		if !suc {
			return
		}

		var num int
		_, err := fmt.Sscanf(m.Content, "donate %d", &num)
		if b.handle(err, m) {
			return
		}

		if user2.Wallet < num {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User <@%s> doesn't even have %d coins!", m.Mentions[0], num))
			return
		}

		if user1.Wallet < num {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("If you are going to rob someone of %d coins, you need to have that many coins in case you get caught.", num))
			return
		}

		num -= rand.Intn(num / 10) // loss

		backNum := rand.Intn(int(math.Pow(float64(ups), 1.5))) - ups - 2 // backfiring
		if backNum < 0 {
			user1.Wallet -= num
			user2.Wallet += num
			b.updateuser(m, user1)
			b.updateuser(m, user2)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Oh no! You got caught and had to give %d coins to the person you were stealing from! Try upgrading property `rob` to reduce the chances of backfiring!", num))
			return
		}
		user1.Wallet += num
		user2.Wallet -= num
		b.updateuser(m, user1)
		b.updateuser(m, user2)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Everything went perfectly and you just stole %d coins!", num))
		return
	}
}