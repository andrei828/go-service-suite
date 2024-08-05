
Use mp4box (build it from source)

Create .mpd file
```bash
MP4Box -dash 2000 -profile dashavc264:live -bs-switching multi -url-template sample.mp4#trackID=1:id=vid0:role=vid0 sample.mp4#trackID=2:id=aud0:role=aud0 -out sample_200.mpd
```

pure luck:
```bash
https://stackoverflow.com/questions/78682030/dash-js-fetches-m4s-files-when-setting-video-currenttime-and-results-in-errors
```