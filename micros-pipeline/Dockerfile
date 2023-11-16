FROM python:3.11

WORKDIR /app

ENV VIRTUAL_ENV=/opt/venv

RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

COPY ./requirements.txt .

RUN pip install --upgrade pip -i https://pypi.tuna.tsinghua.edu.cn/simple && \
    pip install --no-cache-dir -r ./requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

RUN useradd -m admin
USER admin

COPY . .

ENTRYPOINT ["python", "./cmd/app/main.py"]