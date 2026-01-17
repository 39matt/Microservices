using Gateway.Model;
using Gateway.Repository;
using Microsoft.AspNetCore.Mvc;

namespace Gateway.Controllers;
    
[ApiController]
[Route("api/[controller]")]
public class ReadingsController : ControllerBase
{
    private readonly IReadingRepository _repo;

    public ReadingsController(IReadingRepository repo)
    {
        _repo = repo;
    }

    [HttpGet]
    public ActionResult<IEnumerable<Reading>> GetAll()
    {
        return Ok(_repo.GetAllReadings());
    }

    [HttpGet("{id:int}")]
    public ActionResult<Reading> Get(int id)
    {
        return Ok(_repo.GetReading(id));
    }

    [HttpPost]
    public IActionResult Create([FromBody] Reading reading)
    {
        _repo.CreateReading(reading);
        return Created("", reading);
    }

    [HttpPut("{id:int}")]
    public IActionResult Update(int id, [FromBody] Reading reading)
    {
        reading.Id = id;
        _repo.UpdateReading(reading);
        return NoContent();
    }

    [HttpDelete("{id:int}")]
    public IActionResult Delete(int id)
    {
        _repo.DeleteReading(id);
        return NoContent();
    }
    
    [HttpDelete]
    public IActionResult DeleteAll()
    {
        _repo.DeleteAllReadings();
        return NoContent();
    }
}