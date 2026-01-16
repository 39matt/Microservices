using Gateway.Protos;
using Reading = Gateway.Model.Reading;

namespace Gateway.Repository;

public class ReadingRepository : IReadingRepository
{
    public readonly ReadingService.ReadingServiceClient _grpcClient;

    public ReadingRepository(ReadingService.ReadingServiceClient grpcClient)
    {
        _grpcClient = grpcClient;
    }
    
    public IEnumerable<Reading> GetAllReadings()
    {
        var result = _grpcClient.GetAllReadings(new Empty());
        
        return result.Readings.Select(MapProtoToModel);
    }

    public Reading GetReading(int id)
    {
        var result = _grpcClient.GetReading(new GetReadingRequest{Id = id.ToString()});
        
        return MapProtoToModel(result);
    }

    public void CreateReading(Reading reading)
    {
        var request = new CreateReadingRequest { Reading = MapModelToProto(reading) };
        
        _grpcClient.CreateReading(request);
    }

    public void UpdateReading(Reading reading)
    {
        var request = new UpdateReadingRequest() { Reading = MapModelToProto(reading) };
        
        _grpcClient.UpdateReading(request);
        
    }

    public void DeleteReading(int id)
    {
        var request = new RemoveReadingRequest() { Id = id.ToString() };
        
        _grpcClient.RemoveReading(request);
        
    }
    
    private static Reading MapProtoToModel(Gateway.Protos.Reading r) =>
        new Reading
        {
            Id = int.TryParse(r.Id, out var i) ? i : 0,
            Timestamp = r.Timestamp,
            DeviceId = r.DeviceId,
            Co = r.Co,
            Humidity = r.Humidity,
            Light = r.Light,
            Lpg = r.Lpg,
            Motion = r.Motion,
            Smoke = r.Smoke,
            Temperature = r.Temperature
        };

    private static Gateway.Protos.Reading MapModelToProto(Reading r) =>
        new Gateway.Protos.Reading
        {
            Id = r.Id.ToString(),
            Timestamp = r.Timestamp,
            DeviceId = r.DeviceId,
            Co = r.Co,
            Humidity = r.Humidity,
            Light = r.Light,
            Lpg = r.Lpg,
            Motion = r.Motion,
            Smoke = r.Smoke,
            Temperature = r.Temperature
        };
}