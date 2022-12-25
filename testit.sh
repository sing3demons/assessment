#!/bin/sh


docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests


docker-compose -f docker-compose.test.yml down