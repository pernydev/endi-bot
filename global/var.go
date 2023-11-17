package global

import (
	"context"
	"endi/ai"

	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
)

var (
	Discord *discordgo.Session
	RedisC  *redis.Client
	Ctx     = context.Background()
	Model   ai.Model
)
