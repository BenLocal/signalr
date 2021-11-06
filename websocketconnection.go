package signalr

import (
	"bytes"
	"context"
	"fmt"

	"github.com/teivah/onecontext"
	"nhooyr.io/websocket"
)

type webSocketConnection struct {
	ConnectionBase
	conn         *websocket.Conn
	transferMode TransferMode
}

func newWebSocketConnection(parentContext context.Context, requestContext context.Context, connectionID string, conn *websocket.Conn) *webSocketConnection {
	ctx, _ := onecontext.Merge(parentContext, requestContext)
	w := &webSocketConnection{
		conn: conn,
		ConnectionBase: ConnectionBase{
			ctx:          ctx,
			connectionID: connectionID,
		},
	}
	return w
}

func (w *webSocketConnection) Write(p []byte) (n int, err error) {
	if err := w.Context().Err(); err != nil {
		return 0, fmt.Errorf("webSocketConnection canceled: %w", w.ctx.Err())
	}
	ctx := w.ctx
	if w.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(w.ctx, w.Timeout())
		defer cancel() // has no effect because timeoutCtx is either done or not used anymore after websocket returns. But it keeps lint quiet
	}
	messageType := websocket.MessageText
	if w.transferMode == BinaryTransferMode {
		messageType = websocket.MessageBinary
	}
	err = w.conn.Write(ctx, messageType, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w *webSocketConnection) Read(p []byte) (n int, err error) {
	if err := w.Context().Err(); err != nil {
		return 0, fmt.Errorf("webSocketConnection canceled: %w", w.ctx.Err())
	}
	ctx := w.ctx
	if w.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(w.ctx, w.Timeout())
		defer cancel() // has no effect because timeoutCtx is either done or not used anymore after websocket returns. But it keeps lint quiet
	}
	_, data, err := w.conn.Read(ctx)
	if err != nil {
		return 0, err
	}
	return bytes.NewReader(data).Read(p)
}

func (w *webSocketConnection) TransferMode() TransferMode {
	return w.transferMode
}

func (w *webSocketConnection) SetTransferMode(transferMode TransferMode) {
	w.transferMode = transferMode
}
