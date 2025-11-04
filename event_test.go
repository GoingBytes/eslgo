/*
 * Copyright (c) 2020 Percipia
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Contributor(s):
 * Andrew Querol <aquerol@percipia.com>
 */
package eslgo

import (
	"io"
	"net"
	"net/textproto"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestEventToSend = "Content-Length: 483\r\nContent-Type: text/event-plain\r\n\r\nMessage-Account: sip%3A1006%4010.0.1.250\r\nEvent-Name: MESSAGE_QUERY\r\nCore-UUID: 2130a7d1-c1f7-44cd-8fae-8ed5946f3cec\r\nFreeSWITCH-Hostname: localhost.localdomain\r\nFreeSWITCH-IPv4: 10.0.1.250\r\nFreeSWITCH-IPv6: 127.0.0.1\r\nEvent-Date-Local: 2007-12-16%2022%3A29%3A59\r\nEvent-Date-GMT: Mon,%2017%20Dec%202007%2004%3A29%3A59%20GMT\r\nEvent-Date-timestamp: 1197865799573052\r\nEvent-Calling-File: sofia_reg.c\r\nEvent-Calling-Function: sofia_reg_handle_register\r\nEvent-Calling-Line-Number: 603\r\n\r\n"

func TestEvent_readPlainEvent(t *testing.T) {
	server, client := net.Pipe()
	connection := newConnection(client, false, DefaultOptions)
	defer connection.Close()
	defer server.Close()
	defer client.Close()

	var wait sync.WaitGroup
	wait.Add(1)
	connection.RegisterEventListener(EventListenAll, func(event *Event) {
		assert.NotNil(t, event)
		assert.Equal(t, "MESSAGE_QUERY", event.GetName())
		assert.Len(t, event.Headers, 12)
		wait.Done()
	})

	_, err := server.Write([]byte(TestEventToSend))
	assert.Nil(t, err)
	wait.Wait()
}

func TestReadPlainEventWithBody(t *testing.T) {
	event, err := readPlainEvent([]byte("Content-Length: 5\r\nEvent-Name: CUSTOM\r\n\r\nhello"))
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "CUSTOM", event.GetName())
	assert.Equal(t, []byte("hello"), event.Body)
	assert.Len(t, event.Headers, 2)
}

func TestReadPlainEventInvalidContentLength(t *testing.T) {
	event, err := readPlainEvent([]byte("Content-Length: nope\r\n\r\n"))
	assert.Error(t, err)
	assert.NotNil(t, event)
	assert.True(t, event.HasHeader("Content-Length"))
	assert.Nil(t, event.Body)
}

func TestReadPlainEventShortBody(t *testing.T) {
	event, err := readPlainEvent([]byte("Content-Length: 5\r\n\r\nhey"))
	assert.Error(t, err)
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
	assert.NotNil(t, event)
	assert.Len(t, event.Body, 5)
	assert.Equal(t, "hey", string(event.Body[:3]))
}

func TestReadJSONEvent(t *testing.T) {
	payload := []byte(`{"foo":"bar"}`)
	event := readJSONEvent(payload)
	assert.NotNil(t, event)
	assert.Len(t, event.Headers, 0)
	assert.Equal(t, payload, event.Body)
}

func TestEventHeaderHelpers(t *testing.T) {
	event := Event{
		Headers: textproto.MIMEHeader{
			"Event-Name":      []string{"MESSAGE_QUERY"},
			"Message-Account": []string{"sip%3A1006%4010.0.1.250"},
		},
	}

	assert.True(t, event.HasHeader("Event-Name"))
	assert.False(t, event.HasHeader("Missing"))
	assert.Equal(t, "MESSAGE_QUERY", event.GetName())
	assert.Equal(t, "sip:1006@10.0.1.250", event.GetHeader("Message-Account"))
}

func TestEventStringIncludesHeadersAndBody(t *testing.T) {
	event := Event{
		Headers: textproto.MIMEHeader{
			"Event-Name": []string{"TEST_EVENT"},
			"Foo":        []string{"bar"},
		},
		Body: []byte("payload-data"),
	}

	asString := event.String()
	assert.True(t, strings.HasPrefix(asString, "TEST_EVENT\n"))
	assert.Contains(t, asString, "Event-Name:")
	assert.Contains(t, asString, "Foo:")
	assert.Contains(t, asString, "payload-data")
}
