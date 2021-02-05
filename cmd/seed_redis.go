package cmd

import (
	"fmt"
	"os"

	"github.com/DoNewsCode/std/pkg/contract"
	"github.com/DoNewsCode/std/pkg/core"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

// NewSeedRedisCommand returns an command that shows how to retrieve an arbitrary information from
// registered modules
func NewSeedRedisCommand(c *core.C) *cobra.Command {
	var force bool

	type redisSeeder interface {
		SeedRedis() func(client redis.UniversalClient) error
	}
	var seedCmd = &cobra.Command{
		Use:   "seedRedis",
		Short: "seed the redis",
		Long:  `seed the redis by injecting redis instance into the redisSeeder interface. redisSeeder should be registered beforehand.`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = c.Invoke(func(env contract.Env) {
				if env.IsProduction() && !force {
					c.Err(fmt.Errorf("seeding in production requires force flag to be set"))
					os.Exit(1)
				}
			})

			err := c.Modules.Filter(func(seeder redisSeeder) error {
				return c.Invoke(seeder.SeedRedis())
			})
			c.CheckErr(err)

			c.Info("seeding successfully completed")
		},
	}
	seedCmd.Flags().BoolVarP(&force, "force", "f", false, "seeding in production requires force flag to be set")

	return seedCmd
}