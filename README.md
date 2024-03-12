# golangopentelemetry

## 1. Description
해당 코드는 [opentelemtry](https://opentelemetry.io/docs/languages/go/), [smaple code](https://github.com/open-telemetry/opentelemetry-go/blob/main/example/otel-collector/main.go)를 참고하여 작성되었습니다.  
Trace를 사용할 때 필요한 셋팅을 해주는 코드입니다.

## 2. How to use
### 2.1. Setup()
사용할 때 우선 Setup을 먼저 해줘야 합니다.  
Setup은 trace의 기본 설정을 해주는 함수입니다.  

### 2.2. Shutdown()
사용이 끝난 후에는 Shutdown을 해주어야 합니다.  
Trace는 기본적으로 context를 사용하기 때문에, 내부적으로 열려있는 소켓, 고루틴 등을 닫아주어야 합니다.  

### 2.3. GetTracer()
Trace를 사용할 때 span을 직접 넣어주려면 GetTracer를 사용하여 Tracer를 가져와야 합니다.  
사용할때는 우선 Setup()을 해준 후 사용해야 합니다.  
[example](https://github.com/open-telemetry/opentelemetry-go/blob/main/example/otel-collector/main.go)를 참고하여 사용하면 됩니다.
