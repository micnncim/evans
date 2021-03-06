package env_test

import (
	"testing"

	"github.com/ktr0731/evans/entity"
	"github.com/ktr0731/evans/entity/env"
	mockentity "github.com/ktr0731/evans/tests/mock/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	headers := []entity.Header{{Key: "foo", Val: "bar"}}

	t.Run("New", func(t *testing.T) {
		env := env.New(nil, headers)
		h := env.Headers()
		require.Len(t, h, 1)
		require.Equal(t, h[0].Key, "foo")
		require.Equal(t, h[0].Val, "bar")
	})

	t.Run("NewFromServices", func(t *testing.T) {
		env := env.NewFromServices(nil, nil, headers)
		assert.Equal(t, "default", env.DSN())
	})
}

func TestEnv(t *testing.T) {
	pkgs := []*entity.Package{
		{
			Name: "helloworld",
			Services: []entity.Service{
				&mockentity.ServiceMock{
					NameFunc: func() string { return "Greeter" },
					RPCsFunc: func() []entity.RPC {
						return []entity.RPC{
							&mockentity.RPCMock{
								NameFunc: func() string { return "SayHello" },
							},
						}
					},
				},
			},
			Messages: []entity.Message{
				&mockentity.MessageMock{
					FieldsFunc: func() []entity.Field {
						return []entity.Field{
							&mockentity.FieldMock{
								FieldNameFunc: func() string { return "name" },
							},
						}
					},
					NameFunc: func() string { return "HelloRequest" },
				},
				&mockentity.MessageMock{
					NameFunc: func() string { return "HelloResponse" },
				},
			},
		},
	}
	setup := func(t *testing.T) *env.Env {
		return env.New(pkgs, nil)
	}

	env := setup(t)

	t.Run("DSN with no current package", func(t *testing.T) {
		assert.Equal(t, "", env.DSN())
	})

	t.Run("HasCurrentPackage", func(t *testing.T) {
		require.False(t, env.HasCurrentPackage())
		env.UsePackage("helloworld")
		require.True(t, env.HasCurrentPackage())
	})

	t.Run("DSN with no current service", func(t *testing.T) {
		assert.Equal(t, "helloworld", env.DSN())
	})

	t.Run("HasCurrentService", func(t *testing.T) {
		require.False(t, env.HasCurrentService())
		env.UseService("Greeter")
		require.True(t, env.HasCurrentService())
	})

	t.Run("DSN", func(t *testing.T) {
		assert.Equal(t, "helloworld.Greeter", env.DSN())
	})

	t.Run("Packages", func(t *testing.T) {
		pkgs := env.Packages()
		require.Len(t, pkgs, 1)
	})

	t.Run("Services", func(t *testing.T) {
		svcs, err := env.Services()
		require.NoError(t, err)
		require.Len(t, svcs, 1)
	})

	t.Run("Messages", func(t *testing.T) {
		msgs, err := env.Messages()
		require.NoError(t, err)
		require.Len(t, msgs, 2)
	})

	t.Run("RPCs", func(t *testing.T) {
		rpcs, err := env.RPCs()
		require.NoError(t, err)
		require.Len(t, rpcs, 1)
	})

	t.Run("Service", func(t *testing.T) {
		svc, err := env.Service("Greeter")
		require.NoError(t, err)
		require.Equal(t, "Greeter", svc.Name())
		require.Len(t, svc.RPCs(), 1)
	})

	t.Run("Message", func(t *testing.T) {
		msg, err := env.Message("HelloRequest")
		require.NoError(t, err)
		require.Equal(t, "HelloRequest", msg.Name())
		require.Len(t, msg.Fields(), 1)
	})

	t.Run("RPC", func(t *testing.T) {
		rpc, err := env.RPC("SayHello")
		require.NoError(t, err)
		require.Equal(t, "SayHello", rpc.Name())
	})

	t.Run("AddHeader", func(t *testing.T) {
		env := setup(t)
		require.Len(t, env.Headers(), 0)

		env.AddHeader(&entity.Header{"megumi", "kato", false})
		assert.Len(t, env.Headers(), 1)

		env.AddHeader(&entity.Header{"megumi", "kato", false})
		assert.Len(t, env.Headers(), 1)

		assert.EqualValues(t, &entity.Header{Key: "megumi", Val: "kato"}, env.Headers()[0])
	})

	t.Run("RemoveHeader", func(t *testing.T) {
		env := setup(t)
		require.Len(t, env.Headers(), 0)

		env.RemoveHeader("foo")

		headers := map[string]string{
			"hazuki":   "katou",
			"kumiko":   "oumae",
			"reina":    "kousaka",
			"sapphire": "kawashima",
		}
		for k, v := range headers {
			env.AddHeader(&entity.Header{k, v, false})
		}
		assert.Len(t, env.Headers(), 4)

		// Headers must be return slice which ordered by key with ASC
		assert.Equal(t, env.Headers()[0].Key, "hazuki")
		assert.Equal(t, env.Headers()[1].Key, "kumiko")
		assert.Equal(t, env.Headers()[2].Key, "reina")
		assert.Equal(t, env.Headers()[3].Key, "sapphire")

		env.RemoveHeader("foo")
		assert.Len(t, env.Headers(), 4)

		env.RemoveHeader("sapphire")
		assert.Len(t, env.Headers(), 3)
		assert.Equal(t, env.Headers()[2].Key, "reina")

		env.RemoveHeader("hazuki")
		assert.Len(t, env.Headers(), 2)
		assert.Equal(t, env.Headers()[1].Key, "reina")
	})
}
