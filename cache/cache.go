package cache

import (
	pb "RAE_AGE2_JAVAWEB/pdcpp/pdc_trade"
	"log"
	"runtime"
	"sync"
)

var instance *cache
var once sync.Once

//Cache Represents a caching system with pb.PdcTradePrices being the cached items
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

func (c *cache) GetPrices(tID string, gpIDX string) *pb.PdcTradePrices {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return c.activeMap[keyComposer(tID, gpIDX)]
}

func (c *cache) GetPricesCombo(instMap map[string]float64, gpIDX string) *pb.PdcTradePrices {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return &pb.PdcTradePrices{
		TradeId:  "",
		GpIndex:  gpIDX,
		Pv:       sumSingle(instMap, func(s string) float64 { return c.activeMap[s].GetPv() }),
		BondPv:   sumSingle(instMap, func(s string) float64 { return c.activeMap[s].GetBondPv() }),
		TcPrices: comboMultiplePrices(instMap, func(s string) []float64 { return c.activeMap[s].GetTcPrices() }),
		TgPrices: comboMultiplePrices(instMap, func(s string) []float64 { return c.activeMap[s].GetTgPrices() }),
		BondTc:   comboMultiplePrices(instMap, func(s string) []float64 { return c.activeMap[s].GetBondTc() }),
		BondTg:   comboMultiplePrices(instMap, func(s string) []float64 { return c.activeMap[s].GetBondTg() }),
	}
}

func sumSingle(instMap map[string]float64, f func(s string) float64) float64 {
	sum := 0.0
	for k, v := range instMap {
		sum += f(k) * v
	}
	return sum
}

func comboMultiplePrices(instMap map[string]float64, f func(string) []float64) []float64 {

	var sum []float64
	for k, m := range instMap {
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

func (c *cache) SetIdxGpComplete(gpIDX string) {
	gpLock := new(sync.Mutex)
	gpLock.Lock()
	defer gpLock.Unlock()
	log.Printf("Filling finished for gpIdx %v", gpIDX)
	c.gprwLock.Lock()
	c.gpSet[gpIDX] = true
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

func (c *cache) Contains(tID string) bool {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	_, ok := c.activeMap[keyComposer(tID, "0")]
	return ok
}

func (c *cache) GetPv(tID string) float64 {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return c.activeMap[keyComposer(tID, "0")].GetPv()
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

//GetCache returns a cache singleton struct
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
