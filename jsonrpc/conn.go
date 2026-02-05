package jsonrpc

import (
	"bufio"
	"context"
	j "encoding/json"
	"fmt"
	"sync"
)

type Conn struct {
	r         *bufio.Reader
	w         *bufio.Writer
	closeFn   func() error
	incoming  chan Message
	outgoing  chan *Response
	errors    chan error
	once      sync.Once
	ctx       context.Context
	cancel    context.CancelFunc
	writeDone chan struct{}
}

func NewConn(reader *bufio.Reader, writer *bufio.Writer, closeFn func() error) *Conn {
	if closeFn == nil {
		panic("closeFn must be provided to unblock reader")
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Conn{
		r:         reader,
		w:         writer,
		closeFn:   closeFn,
		incoming:  make(chan Message, 100),
		outgoing:  make(chan *Response, 100),
		errors:    make(chan error, 2),
		ctx:       ctx,
		cancel:    cancel,
		writeDone: make(chan struct{}),
	}
}

func (c *Conn) Write(resp *Response) error {
	content, err := j.Marshal(resp)
	if err != nil {
		return err
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(content))

	if _, err := c.w.WriteString(header); err != nil {
		return err
	}

	if _, err := c.w.Write(content); err != nil {
		return err
	}

	return c.w.Flush()
}

func (c *Conn) readLoop() {
	defer func() {
		if r := recover(); r != nil {
			c.fail(fmt.Errorf("panic in readLoop: %v", r))
		}
		close(c.incoming)
	}()
	for {
		msg, err := readBody(c.r)
		if err != nil {
			c.fail(fmt.Errorf("read error: %w", err))
			return
		}
		select {
		case c.incoming <- msg:
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Conn) writeLoop() {
	defer close(c.writeDone)
	for resp := range c.outgoing {
		if err := c.Write(resp); err != nil {
			c.fail(fmt.Errorf("write error: %w", err))
			return
		}
	}
}

func (c *Conn) SendResponse(resp *Response) error {
	select {
	case <-c.ctx.Done():
		return fmt.Errorf("connection closed")
	default:
	}

	select {
	case c.outgoing <- resp:
		return nil

	case <-c.ctx.Done():
		return fmt.Errorf("connection closed")
	}
}

func (c *Conn) Incoming() <-chan Message {
	return c.incoming
}

func (c *Conn) Close() error {
	c.fail(nil)
	c.Wait()
	return nil
}

func (c *Conn) Start() {
	go c.readLoop()
	go c.writeLoop()
}

func (c *Conn) reportError(err error) {
	if err == nil {
		return
	}
	select {
	case c.errors <- err:
	default:
	}
}

func (c *Conn) fail(err error) {
	c.once.Do(func() {
		c.reportError(err)

		close(c.outgoing)

		_ = c.closeFn()

		c.cancel()
		close(c.errors)
	})
}

// Wait blocks until shutdown is initiated (fail called). It does not guarantee readLoop exits
// unless closeFn unblocks the reader.
func (c *Conn) Wait() {
	<-c.ctx.Done()
	<-c.writeDone
}
