using Gateway.Model;

namespace Gateway.Repository;

public interface IReadingRepository
{
    IEnumerable<Reading> GetAllReadings();
    Reading GetReading(int id);
    void CreateReading(Reading reading);
    void UpdateReading(Reading reading);
    void DeleteReading(int id);
    void Save();
}