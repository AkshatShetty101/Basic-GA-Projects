package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jinzhu/copier"
)

type Resource struct {
	name   string
	value  float64
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

type Chromosome struct {
	genes []itemQuantity
	f     fitness
	age   int
}

func main() {
	start := time.Now()
	items := []Resource{Resource{"Flour", 1680.0, 0.265, 0.41}, Resource{"Butter", 1440.0, 0.5, 0.13}, Resource{"Sugar", 1840.0, 0.441, .29}}
	parent := Chromosome{}
	bestParent := Chromosome{}
	child := Chromosome{}
	optimal := getFitness([]itemQuantity{itemQuantity{items[0], 1}, itemQuantity{items[1], 14}, itemQuantity{items[2], 6}})
	fmt.Println("\nOptimal Fitness:")
	fmt.Printf("%+v\n", optimal)
	maxWeight := 10.0
	maxVolume := 4.0
	parent.genes = create(parent.genes, items, maxWeight, maxVolume)
	parent.f = getFitness(parent.genes)
	parent.age = 0
	copier.Copy(&bestParent, &parent)
	display(parent, maxWeight, maxVolume, start)
	for {
		child.genes = mutate(parent.genes, items, maxWeight, maxVolume)
		child.f = getFitness(child.genes)
	}
}

func compareFitness(c fitness, p fitness) bool {
	return c.totalValue > p.totalValue
}

func display(ob Chromosome, maxW float64, maxV float64, start time.Time) {
	fmt.Println("-----------------------------------------------")
	fmt.Println("Knapsack Details:")
	fmt.Println("Name\tValue\t\tWeight\tVolume\tQuantity")
	for _, val := range ob.genes {
		fmt.Printf("%s\t%0.2f\t\t%0.2f\t%0.2f\t%v\n\n", val.item.name, val.item.value, val.item.weight, val.item.volume, val.quantity)
	}
	fmt.Println("Fitness Details:")
	fmt.Printf("Total Value: %0.2f\nTotal Weight: %0.2f\t\tMaximum Weight: %0.2f\nTotal Volume: %0.2f\t\tMaximum Volume: %0.2f\n", ob.f.totalValue, ob.f.totalWeight, maxW, ob.f.totalVolume, maxV)
}

func getFitness(iq []itemQuantity) fitness {
	totalWeight := 0.0
	totalVolume := 0.0
	totalValue := 0.0
	for _, item := range iq {
		totalValue += item.item.value * float64(item.quantity)
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
	if maxQ > 0.0 {
		return itemQuantity{items[x], maxQ}
	}
	return itemQuantity{}
}

func mutate(kn []itemQuantity, items []Resource, maxW float64, maxV float64) []itemQuantity {
	k := []itemQuantity{}
	copier.Copy(&k, &kn)
	Fit := getFitness(k)
	remainingW := maxW - Fit.totalWeight
	remainingV := maxV - Fit.totalVolume
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	//Removing from the knapsack
	if len(k) > 1 && r.Intn(10) == 1 {
		// fmt.Println("########### Removing Item ##################")
		// fmt.Println()
		idx := r.Intn(len(k))
		remainingW += k[idx].item.weight * float64(k[idx].quantity)
		remainingV += k[idx].item.volume * float64(k[idx].quantity)
		k[len(k)-1], k[idx] = k[idx], k[len(k)-1]
		k = k[:len(k)-1]
	}
	//Adding to the knapsack
	if ((remainingW > 0.0 || remainingV > 0.0) && len(k) == 0) || (len(k) < len(items) && r.Intn(10) == 1) {
		// fmt.Println("########### Adding New Item ##################")
		newItem := add(k, items, remainingW, remainingV)
		if newItem != (itemQuantity{}) {
			k = append(k, newItem)
			remainingW -= newItem.item.weight * float64(newItem.quantity)
			remainingV -= newItem.item.volume * float64(newItem.quantity)
		}
	}
	//Repalcing item in the knapsack
	// fmt.Println("########### Replacing Item ##################")
	idx := r.Intn(len(k))
	// fmt.Println(k[idx].item.weight, float64(k[idx].quantity))
	remainingW += k[idx].item.weight * float64(k[idx].quantity)
	remainingV += k[idx].item.volume * float64(k[idx].quantity)
	// fmt.Println(remainingW, remainingV)
	list := r.Perm(len(items))
	itema := items[list[0]]
	itemb := items[list[1]]
	// fmt.Println(itema, itemb, k[idx].item)
	if itema != k[idx].item {
		maxQ := maxQuantity(itema, remainingW, remainingV)
		// fmt.Println("\n\n\n\n", maxQ)
		if maxQ > 0.0 {
			k[idx] = itemQuantity{itema, maxQ}
		} else {
			k[len(k)-1], k[idx] = k[idx], k[len(k)-1]
			k = k[:len(k)-1]
		}

	} else {
		maxQ := maxQuantity(itemb, remainingW, remainingV)
		// fmt.Println("\n\n\n\n", maxQ)
		if maxQ > 0.0 {
			k[idx] = itemQuantity{itemb, maxQ}
		} else {
			k[len(k)-1], k[idx] = k[idx], k[len(k)-1]
			k = k[:len(k)-1]
		}
	}
	return k
}
