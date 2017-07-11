
## concourwe-webpage-resources

See `sample_pipeline.yml` for an example.

## Development

```bash
ginkgo -r
make
docker build . -t cholick/concourse-webpage-resource
docker push cholick/concourse-webpage-resource
```
