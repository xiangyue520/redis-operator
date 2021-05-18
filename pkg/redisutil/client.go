package redisutil

import (
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// IClient redis client interface
type IClient interface {
	// Close closes the connection.
	Close() error

	// Cmd calls the given Redis command.
	Cmd(cmd string, args ...interface{}) *redis.Cmd

	// PipeAppend adds the given call to the pipeline queue.
	// Use PipeResp() to read the response.
	PipeAppend(cmd string, args ...interface{})

	// PipeResp returns the reply for the next request in the pipeline queue. Err
	// with ErrPipelineEmpty is returned if the pipeline queue is empty.
	PipeResp() *redis.Cmd

	// PipeClear clears the contents of the current pipeline queue, both commands
	// queued by PipeAppend which have yet to be sent and responses which have yet
	// to be retrieved through PipeResp. The first returned int will be the number
	// of pending commands dropped, the second will be the number of pending
	// responses dropped
	PipeClear() (int, int)

	// ReadResp will read a Resp off of the connection without sending anything
	// first (useful after you've sent a SUSBSCRIBE command). This will block until
	// a reply is received or the timeout is reached (returning the IOErr). You can
	// use IsTimeout to check if the Resp is due to a Timeout
	//
	// Note: this is a more low-level function, you really shouldn't have to
	// actually use it unless you're writing your own pub/sub code
	//ReadResp() *redis.Cmd
	Client() *redis.Client
}

// Client structure representing a client connection to redis
type Client struct {
	commandsMapping map[string]string
	client          *redis.Client
}

// NewClient build a client connection and connect to a redis address
func NewClient(addr, password string, cnxTimeout time.Duration, commandsMapping map[string]string) (IClient, error) {
	var err error
	c := &Client{
		commandsMapping: commandsMapping,
	}

	//c.client, err = redis.DialTimeout("tcp", addr, cnxTimeout)
	c.client = redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: 0,
		DialTimeout: cnxTimeout,
	})
	//执行一次ping命令,然后进行处理
	ctx := c.client.Context()
	//cmd:=c.client.Ping(ctx)
	//fmt.Println(cmd.Val())

	if password != "" {
		err= c.client.Do(ctx,"AUTH", password).Err()
		//err = c.client.Cmd("AUTH", password).Err
	}
	return c, err
}

// Close closes the connection.
func (c *Client) Close() error {
	return c.client.Close()
}

// Cmd calls the given Redis command.
func (c *Client) Cmd(cmd string, args ...interface{}) *redis.Cmd {

	return c.client.Do(c.client.Context(),c.getCommand(cmd), args)
}

// getCommand return the command name after applying rename-command
func (c *Client) getCommand(cmd string) string {
	upperCmd := strings.ToUpper(cmd)
	if renamed, found := c.commandsMapping[upperCmd]; found {
		return renamed
	}
	return upperCmd
}

// PipeAppend adds the given call to the pipeline queue.
func (c *Client) PipeAppend(cmd string, args ...interface{}) {
	ctx:=c.client.Context()
	//c.client.Pipeline().Append(ctx,c.getCommand(cmd),args)
	c.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		return pipe.Do(ctx,c.getCommand(cmd),args).Err()
	})
}

// PipeResp returns the reply for the next request in the pipeline queue. Err
func (c *Client) PipeResp() *redis.Cmd {
	ctx:=c.client.Context()
	//todo 需要做pipeline的回复
	return c.client.Pipeline().Do(ctx,"ping")
}

// PipeClear clears the contents of the current pipeline queue
func (c *Client) PipeClear() (int, int) {
	return 1,2
}

func (c *Client)Client() *redis.Client{
	return c.client
}

// ReadResp will read a Resp off of the connection without sending anything
//func (c *Client) ReadResp() *redis.Cmd {
//	return c.client.Read
//}
