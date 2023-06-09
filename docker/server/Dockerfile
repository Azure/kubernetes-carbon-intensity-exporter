FROM mcr.microsoft.com/dotnet/sdk:6.0 AS build-env
WORKDIR /app

# Copy everything from source
RUN git clone https://github.com/helayoty/carbon-aware-sdk.git && cd carbon-aware-sdk && git checkout 55d75c7
# Use implicit restore to build and publish
RUN dotnet publish carbon-aware-sdk/src/CarbonAware.WebApi/src/CarbonAware.WebApi.csproj -c Release -o publish
# Generate OpenAPI spec
WORKDIR carbon-aware-sdk/src/CarbonAware.WebApi/src
RUN dotnet tool restore && \
    mkdir -p /app/publish/wwwroot/api/v1 && \
    dotnet tool run swagger tofile --output /app/publish/wwwroot/api/v1/swagger.yaml --yaml /app/publish/CarbonAware.WebApi.dll v1

# Build runtime image
FROM mcr.microsoft.com/dotnet/aspnet:6.0
WORKDIR /app
COPY --from=build-env /app/publish .
ENTRYPOINT ["dotnet", "CarbonAware.WebApi.dll"]
