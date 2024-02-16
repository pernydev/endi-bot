package global

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
	"github.com/supabase-community/supabase-go"
)

var (
	Discord   *discordgo.Session
	RedisC    *redis.Client
	Ctx       = context.Background()
	SupabaseC *supabase.Client
)
