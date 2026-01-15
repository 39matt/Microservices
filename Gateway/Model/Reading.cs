namespace Gateway.Model;

public class Reading
{
    public int Id { get; set; }
    public string Timestamp { get; set; }
    public string DeviceId { get; set; }
    public double Co { get; set; }
    public float Humidity { get; set; }
    public bool Light { get; set; }
    public double Lpg { get; set; }
    public bool Motion { get; set; }
    public double Smoke { get; set; }
    public float Temperature { get; set; }
}