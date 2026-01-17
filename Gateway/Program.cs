using Gateway.Protos;
using Gateway.Repository;
using Grpc.Core;
using Grpc.Net.Client;

var builder = WebApplication.CreateBuilder(args);

// builder.Services.AddOpenApi();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
builder.Services.AddControllers();
builder.Services.AddGrpcClient<ReadingService.ReadingServiceClient>(o =>
{
    var grpcHost = Environment.GetEnvironmentVariable("DATAMANAGER_HOST") ?? "localhost:8080";
    o.Address = new Uri($"http://{grpcHost}");
});

builder.Services.AddScoped<IReadingRepository, ReadingRepository>();

var app = builder.Build();

// if (app.Environment.IsDevelopment())
// {
    // app.MapOpenApi();
    app.UseSwagger();
    app.UseSwaggerUI(options =>
    {
        options.SwaggerEndpoint("/swagger/v1/swagger.json", "Gateway API");
        options.RoutePrefix = string.Empty;
    });
// }

app.UseHttpsRedirection();
app.MapControllers();

// using var channel = GrpcChannel.ForAddress("http://localhost:8080");
// var client = new ReadingService.ReadingServiceClient(channel);
//
// var reply = await client.GetAllReadingsAsync(new Empty(), new CallOptions());
// Console.WriteLine(reply.Readings);

app.Run();


