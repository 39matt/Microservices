using Gateway.Model;

namespace Gateway.Repository;

public interface IReadingRepository
{
    IEnumerable<Reading> GetReadings();
    Reading GetReadingById(int id);
    void AddReading(Reading reading);
    void UpdateReading(Reading reading);
    void DeleteReading(int id);
    void Save();
}