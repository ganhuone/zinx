package main

import (
	"fmt"
	"io"
	"net"
	"testing"
	"zinx/znet"
)

func TestLen(T *testing.T) {
	m := []byte{'z', 'z'}

	fmt.Println(len(m))
	fmt.Println(m)

}

func TestDataPack(T *testing.T) {
	listnner, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			conn, err := listnner.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}

			go func(conn net.Conn) {

				for {

					dp := znet.NewDataPack()
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println(err)
						break
					}

					msg, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println(err)
						break
					}

					if msg.GetDataLen() > 0 {
						data := make([]byte, msg.GetDataLen())

						_, err := io.ReadFull(conn, data)
						if err != nil {
							fmt.Println(err)
							return
						}

						msg.SetData(data)
					}

					fmt.Println("--Recv MsgID:", msg.GetMsgId(), ", dataLen = ", msg.GetDataLen(), ", data = ", string(msg.GetData()))

				}

			}(conn)

		}
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
	}

	dp := znet.NewDataPack()

	msgOne := &znet.Message{
		Id:   1,
		Data: []byte{'a', 'b', 'c', 'd'},
	}

	msgOne.SetDataLen(uint32(len(msgOne.GetData())))

	sendDataOne, err := dp.Pack(msgOne)
	if err != nil {
		fmt.Println(err)
		return
	}

	msgTwo := &znet.Message{
		Id:   1,
		Data: []byte{'h', 'e', 'l', 'l', 'o'},
	}

	msgTwo.SetDataLen(uint32(len(msgTwo.GetData())))

	sendDataTwo, err := dp.Pack(msgTwo)
	if err != nil {
		fmt.Println(err)
		return
	}

	sendDataOne = append(sendDataOne, sendDataTwo...)

	conn.Write(sendDataOne)

	select {}

}
