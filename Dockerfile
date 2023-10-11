FROM 699983977898.dkr.ecr.us-east-1.amazonaws.com/datadog-agent:boh-dan-datadog-agent-3c420f2-bohdan-7.29.0-with-high-cpu-fix
RUN sed -i -re 's/([a-z]{2}.)?archive.ubuntu.com|security.ubuntu.com/old-releases.ubuntu.com/g' /etc/apt/sources.list
RUN apt-get update && apt-get install -y conntrack
