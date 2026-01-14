namespace Gateway.Model;

public class Reading
{
    public int Id { get; set; }
    public string Timestamp { get; set; }
    public string DeviceId { get; set; }
    public string Co { get; set; }
    public string Humidity { get; set; }
    public string Light { get; set; }
    public string LPG { get; set; }
    public string Motion { get; set; }
    public string Smoke { get; set; }
    public string Temperature { get; set; }
}