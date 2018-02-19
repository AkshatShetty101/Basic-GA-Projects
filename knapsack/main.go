package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Resource struct {
	name   string
	value  int
	weight float64
	volume float64
}

type itemQuantity struct {
	item     Resource
	quantity int
}

type fitness struct {
	totalWeight float64
	totalVolume float64
	totalValue  float64
}

func main() {
	items := []Resource{Resource{"Flour", 1680, 0.265, 0.41}, Resource{"Butter", 1440, 0.5, 0.13}, Resource{"Sugar", 1840, 0.441, .29}}
	knapsack := []itemQuantity{}
	maxWeight := 10.0
	maxVolume := 4.0
	fmt.Println(create(knapsack, items, maxWeight, maxVolume))
	fmt.Println(maxVolume, maxWeight, items)
}

func getFitness(iq []itemQuantity) fitness {
	totalWeight := 0.0
	totalVolume := 0.0
	totalValue := 0.0
	for _, item := range iq {
		totalValue += item.item.value * item.quantity
		totalWeight += item.item.weight * float64(item.quantity)
		totalVolume += item.item.volume * float64(item.quantity)
	}
	return fitness{totalWeight, totalVolume, totalValue}
}

func maxQuantity(item Resource, maxW float64, maxV float64) int {
	return int(math.Min(math.Floor(maxW/item.weight), math.Floor(maxV/item.volume)))
}

func create(k []itemQuantity, items []Resource, maxW float64, maxV float64) []itemQuantity {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	remainingV, remainingW := maxV, maxW
	for i := 0; i <= r.Intn(len(items)); i++ {
		newItem := add(k, items, remainingW, remainingV)
		if newItem != (itemQuantity{}) {
			k = append(k, newItem)
			remainingW -= newItem.item.weight * float64(newItem.quantity)
			remainingV -= newItem.item.volume * float64(newItem.quantity)
		}
	}
	return k
}

func add(k []itemQuantity, items []Resource, maxW float64, maxV float64) itemQuantity {
	usedItems := []Resource{}
	x := 0
	for _, val := range k {
		usedItems = append(usedItems, val.item)
	}
	for {
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		x = r.Intn(len(items))
		flag := 0
		for _, val := range usedItems {
			if items[x] == val {
				flag = -1
			}
		}
		if flag == 0 {
			break
		}
	}
	maxQ := maxQuantity(items[x], maxW, maxV)
	if maxQ > 0 {
		return itemQuantity{items[x], maxQ}
	}
	return itemQuantity{}
}

func mutate(k []itemQuantity, items []Resource, maxW float64, maxV float64) {
	Fit := getFitness(k)
	remainingW := maxW - Fit.totalWeight
	remainingV := maxV - Fit.totalVolume
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	if len(k) > 1 && r.Intn(10) == 0 {
		idx := r.Intn(len(k))
		k[len(k)-1], k[idx] = k[idx], k[len(k)-1]
		k = k[:len(k)-1]
	}
	if ((remainingW > 0 || remainingV > 0) && len(k) == 0) || (len(k) < len(items) && r.Intn(10) == 0) {
		newItem := add(k, items, remainingW, remainingV)
		if newItem != (itemQuantity{}) {
			k = append(k, newItem)
			// remainingW -= newItem.item.weight * float64(newItem.quantity)
			// remainingV -= newItem.item.volume * float64(newItem.quantity)
		}
	}
}
