package main

import (
	"container/heap"
	"errors"
	"math/rand"
)

type stockItem struct {
	averagePrice float64
	buyOrders    PriorityQueue
	sellOrders   PriorityQueue
}

type Exchange struct {
	pq     map[string]*stockItem
	prices map[string]float64
}

func (e *Exchange) describeStock(stock string) (float64, error) {
	if val, ok := e.prices[stock]; ok {
		return val, nil
	}
	return 0.0, errNotFound
}

func NewStockItem() *stockItem {
	buyOrders := make(PriorityQueue, 0)
	heap.Init(&buyOrders)
	sellOrders := make(PriorityQueue, 0)
	heap.Init(&sellOrders)
	return &stockItem{
		averagePrice: 0.0,
		buyOrders:    buyOrders,
		sellOrders:   sellOrders,
	}
}

func (e *Exchange) newStockType(stock string) {
	if _, ok := e.pq[stock]; !ok {
		e.pq[stock] = NewStockItem()
	}
}

var errNotFound = errors.New("could not find stock")

func (e *Exchange) peekBuy(stock string) (*Item, error) {
	if val, ok := e.pq[stock]; ok {
		return heap.Pop(&val.buyOrders).(*Item), nil
	}
	return nil, errNotFound
}

func (e *Exchange) peekSell(stock string) (*Item, error) {
	if val, ok := e.pq[stock]; ok {
		return heap.Pop(&val.sellOrders).(*Item), nil
	}
	return nil, errNotFound
}

func (e *Exchange) addBuy(stock string, item *Item) error {
	if val, ok := e.pq[stock]; ok {
		heap.Push(&val.buyOrders, item)
		return nil
	}
	return errNotFound
}

func (e *Exchange) addSell(stock string, item *Item) error {
	if val, ok := e.pq[stock]; ok {
		heap.Push(&val.sellOrders, item)
		return nil
	}
	return errNotFound
}

type TransactionStatus string

const (
	SUCCESS TransactionStatus = "Success"
	PENDING TransactionStatus = "Pending"
	FAILURE TransactionStatus = "Failure"
)

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func (e *Exchange) buy(stock string, price float64, volume uint64) (TransactionStatus, error) {
	err := e.addBuy(stock, &Item{
		id:     uint64(rand.Int()),
		price:  price,
		volume: volume,
	})
	if err != nil {
		return FAILURE, err
	}

	err = e.matchOrders(stock)
	if err != nil {
		return FAILURE, err
	}

	// hack - we assume it got matched for demo purposes

	return SUCCESS, nil
}

func (e *Exchange) matchOrders(stock string) error {
	for {
		sellOrder, err := e.peekSell(stock)
		if err != nil {
			return err
		}

		if sellOrder == nil {
			return nil
		}

		buyOrder, err := e.peekBuy(stock)
		if err != nil {
			return err
		}

		if buyOrder == nil {
			return nil
		}

		if buyOrder.price >= sellOrder.price {
			amount := min(buyOrder.volume, sellOrder.volume)
			sellOrder.volume -= amount
			buyOrder.volume -= amount
			e.prices[stock] = (buyOrder.price + sellOrder.price) / 2
			if sellOrder.volume != 0 {
				err := e.addSell(stock, sellOrder)
				if err != nil {
					return errNotFound
				}
			}

			if buyOrder.volume != 0 {
				err := e.addBuy(stock, buyOrder)
				if err != nil {
					return errNotFound
				}
			}
		} else {
			return nil
		}
	}
}

func NewExchange() *Exchange {
	return &Exchange{
		pq: make(map[string]*stockItem, 0),
	}
}
