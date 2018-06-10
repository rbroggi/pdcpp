package cache

import (
	pb "RAE_AGE2_JAVAWEB/pdcpp/pdc_trade"
	"log"
	"runtime"
	"sync"
)

var instance *cache
var once sync.Once

//Represents a caching system with pb.PdcTradePrices being the cached items
type Cache interface {
	GetPrices(string, string) *pb.PdcTradePrices
	GetPricesCombo(map[string]float64, string) *pb.PdcTradePrices
	InsertRecord(*pb.PdcTradePrices) bool
	GetCompletedGps() int
	SetIdxGpComplete(string)
	Contains(string) bool
	GetPv(string) float64
	GetActiveContentLen() int
	GetFillingContentLen() int
}

type cache struct {
	fillingMap map[string]*pb.PdcTradePrices
	activeMap  map[string]*pb.PdcTradePrices
	gpSet      map[string]bool
	rwLock     sync.RWMutex
	frwLock    sync.RWMutex
	gprwLock   sync.RWMutex
	gpSize     int
}

func (c *cache) GetPrices(t_id string, gp_idx string) *pb.PdcTradePrices {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return c.activeMap[keyComposer(t_id, gp_idx)]
}

func (c *cache) GetPricesCombo(inst_map map[string]float64, gp_idx string) *pb.PdcTradePrices {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return &pb.PdcTradePrices{
		TradeId:  "",
		GpIndex:  gp_idx,
		Pv:       sumSingle(inst_map, func(s string) float64 { return c.activeMap[s].GetPv() }),
		BondPv:   sumSingle(inst_map, func(s string) float64 { return c.activeMap[s].GetBondPv() }),
		TcPrices: comboMultiplePrices(inst_map, func(s string) []float64 { return c.activeMap[s].GetTcPrices() }),
		TgPrices: comboMultiplePrices(inst_map, func(s string) []float64 { return c.activeMap[s].GetTgPrices() }),
		BondTc:   comboMultiplePrices(inst_map, func(s string) []float64 { return c.activeMap[s].GetBondTc() }),
		BondTg:   comboMultiplePrices(inst_map, func(s string) []float64 { return c.activeMap[s].GetBondTg() }),
	}
}

func sumSingle(inst_map map[string]float64, f func(s string) float64) float64 {
	sum := 0.0
	for k, v := range inst_map {
		sum += f(k) * v
	}
	return sum
}

func comboMultiplePrices(inst_map map[string]float64, f func(string) []float64) []float64 {

	var sum []float64
	for k, m := range inst_map {
		if len(sum) == 0 {
			sum = make([]float64, len(f(k)))
		}
		for i, v := range f(k) {
			sum[i] += v * m
		}
	}
	return sum
}

func (c *cache) InsertRecord(record *pb.PdcTradePrices) bool {
	if record == nil {
		return false
	}
	c.frwLock.Lock()
	defer c.frwLock.Unlock()
	c.fillingMap[keyComposer(record.GetTradeId(), record.GetGpIndex())] = record
	return true
}

func keyComposer(prefix string, suffix string) string {
	return prefix + suffix
}

func (c *cache) GetCompletedGps() int {
	c.gprwLock.RLock()
	defer c.gprwLock.RUnlock()
	return len(c.gpSet)
}

func (c *cache) SetIdxGpComplete(gp_idx string) {
	gpLock := new(sync.Mutex)
	gpLock.Lock()
	defer gpLock.Unlock()
	log.Printf("Filling finished for gpIdx %v", gp_idx)
	c.gprwLock.Lock()
	c.gpSet[gp_idx] = true
	c.gprwLock.Unlock()
	if c.GetCompletedGps() == c.gpSize {
		log.Println("Maps switching...")
		c.rwLock.Lock()
		c.activeMap = c.fillingMap
		c.rwLock.Unlock()
		c.frwLock.Lock()
		c.fillingMap = make(map[string]*pb.PdcTradePrices)
		c.frwLock.Unlock()
		c.gprwLock.Lock()
		c.gpSet = make(map[string]bool)
		c.gprwLock.Unlock()
		log.Println("Maps switched!")
		log.Println("Calling GC...")
		runtime.GC()
		log.Println("GC finished!")

	}
}

func (c *cache) Contains(t_id string) bool {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	_, ok := c.activeMap[keyComposer(t_id, "0")]
	return ok
}

func (c *cache) GetPv(t_id string) float64 {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return c.activeMap[keyComposer(t_id, "0")].GetPv()
}

func (c *cache) GetActiveContentLen() int {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return len(c.activeMap)
}

func (c *cache) GetFillingContentLen() int {
	c.frwLock.RLock()
	defer c.frwLock.RUnlock()
	return len(c.fillingMap)
}

//return cache singleton struct
func GetCache(gpSize int) Cache {
	once.Do(func() {
		instance = &cache{
			fillingMap: make(map[string]*pb.PdcTradePrices),
			activeMap:  make(map[string]*pb.PdcTradePrices),
			gpSet:      make(map[string]bool),
			rwLock:     sync.RWMutex{},
			frwLock:    sync.RWMutex{},
			gprwLock:   sync.RWMutex{},
			gpSize:     gpSize,
		}
	})
	return instance
}
