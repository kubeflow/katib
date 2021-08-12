#This container contains your model and any helper scripts specific to your model.
FROM python:3.7-stretch

# copy requirements
ADD requirements.txt /opt/requirements.txt
# install requirements
ADD requirements.txt /opt/requirements.txt
RUN pip install -r /opt/requirements.txt

# copy model
ADD model.py /opt/model.py
RUN chmod +x /opt/model.py

ENTRYPOINT ["python"]
CMD ["/opt/model.py"]