using Gateway.Protos;
using Microsoft.AspNetCore.Mvc;

[ApiController]
[Route("api/[controller]")]
public class ReadingsController : ControllerBase
{
    private readonly ReadingService.ReadingServiceClient _client;

    public ReadingsController(ReadingService.ReadingServiceClient client)
    {
        _client = client;
    }

    [HttpGet]
    public async Task<IActionResult> GetAll()
    {
        var response = await _client.GetAllReadingsAsync(new Empty());
        return Ok(response.Readings);
    }
}