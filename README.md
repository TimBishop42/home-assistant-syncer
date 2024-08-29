# home-assistant-syncer


# Docker build and deploy
- docker build -t homeassistant/syncer -f docker/Dockerfile . 
- docker tag homeassistant/syncer tbished/syncer:v1
- docker push tbished/syncer:v1