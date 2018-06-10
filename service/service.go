package service

import (
	"RAE_AGE2_JAVAWEB/pdcpp/cache"
	pb "RAE_AGE2_JAVAWEB/pdcpp/pdc_trade"
	context "golang.org/x/net/context"
	"io"
)

func GetPriceProviderService(gpIndex int) pb.PdcTradePricesServiceServer {
	return &pdcPriceProvider{cCache: cache.GetCache(gpIndex)}
}

//service struct
type pdcPriceProvider struct {
	cCache cache.Cache
}

//sync rpc interface to insert records into the cache
func (s *pdcPriceProvider) SendTradePricesRecord(c context.Context, tp *pb.PdcTradePrices) (*pb.Bool, error) {
	return &pb.Bool{Value: s.cCache.InsertRecord(tp)}, nil
}

func (s *pdcPriceProvider) SendTradePricesRecordStream(price_stream pb.PdcTradePricesService_SendTradePricesRecordStreamServer) error {
	var incoming_records int32
	incoming_records = 0
	for {
		record, err := price_stream.Recv()
		if err == io.EOF {
			return price_stream.SendAndClose(&pb.EventNum{Value: incoming_records})
		}
		if err != nil {
			return err
		}
		incoming_records++
		s.cCache.InsertRecord(record)
	}
}

//sync rpc interfacd to fetch a record from the cache
func (s *pdcPriceProvider) GetIdcTradePrices(c context.Context, record *pb.RecordId) (*pb.PdcTradePrices, error) {
	return s.cCache.GetPrices(record.GetTradeId(), record.GetGpIndex()), nil
}

//sync rpc interface to fetch the aggregated resulg of multiple records form the cache
func (s *pdcPriceProvider) GetAggrIdcTradePrices(c context.Context, trade_id_multiplier *pb.TradeIdMultiplierCollection) (*pb.PdcTradePrices, error) {
	return s.cCache.GetPricesCombo(trade_id_multiplier.GetTradeIdMultiplier(), trade_id_multiplier.GetGpIndex()), nil
}

//sync rpc interface to notify the end of a given gpIndex point
func (s *pdcPriceProvider) NotifyEndGpRecords(c context.Context, gpIdx *pb.GpId) (*pb.Bool, error) {
	s.cCache.SetIdxGpComplete(gpIdx.GetGpIndex())
	return &pb.Bool{Value: true}, nil
}

func (s *pdcPriceProvider) Ping(c context.Context, empty *pb.Empty) (*pb.Bool, error) {
	return &pb.Bool{Value: true}, nil
}

func (s *pdcPriceProvider) CheckTrades(c context.Context, trade_ids *pb.TradeIds) (*pb.Bool, error) {
	for _, t := range trade_ids.GetTradeIds() {
		if contains := s.cCache.Contains(t); !contains {
			return &pb.Bool{Value: false}, nil
		}
	}
	return &pb.Bool{Value: true}, nil
}

func (s *pdcPriceProvider) CheckTrade(c context.Context, trade_id *pb.TradeId) (*pb.Bool, error) {
	return &pb.Bool{Value: s.cCache.Contains(trade_id.GetTradeId())}, nil
}

func (s *pdcPriceProvider) GetPv(c context.Context, trade_id *pb.TradeId) (*pb.Pv, error) {
	return &pb.Pv{Pv: s.cCache.GetPv(trade_id.GetTradeId())}, nil
}
