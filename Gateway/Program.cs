using Gateway.Protos;
using Grpc.Core;
using Grpc.Net.Client;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddOpenApi();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
builder.Services.AddControllers();

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
    app.UseSwagger();
    app.UseSwaggerUI(options =>
    {
        options.SwaggerEndpoint("/swagger/v1/swagger.json", "Gateway API");
        options.RoutePrefix = string.Empty;
    });
}

app.UseHttpsRedirection();
app.MapControllers();

using var channel = GrpcChannel.ForAddress("http://localhost:8080");
var client = new ReadingService.ReadingServiceClient(channel);

var reply = await client.GetAllReadingsAsync(new Empty(), new CallOptions());

Console.WriteLine(reply.Readings);
app.Run();


