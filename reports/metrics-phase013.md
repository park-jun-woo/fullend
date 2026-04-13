# Structural Metrics Report (Phase011)

2026-04-13 · pkg/generate vs internal/gen

## 파일·줄 수
                    internal/gen  pkg/generate
파일 수             224           462         
평균 줄 수          34.3          28.1        
중앙값 줄 수        28            23          
최대 줄 수          217           217         

## 함수 매개변수 분포
                    internal/gen  pkg/generate
총 함수 수          208           422         
평균 매개변수       2.56          2.23        
중앙값              2             2           
최대                10            10          
8+ params 함수      12            10          
5+ params 함수      29            29          

## 중복 패턴 *WithDomains
                    internal/gen  pkg/generate
*WithDomains 함수   4             0           

## Decide* 순수 판정 함수 (Phase010 구조 정비 지표)
                    internal/gen  pkg/generate
Decide* 함수 수     0             3           

## Toulmin 사용 (참고)
                    internal/gen  pkg/generate
toulmin.NewGraph    0             0           

> Phase010 결정: 2-depth 이내 if-else 로 해결되어 Toulmin 미채택. 대신 Decide* 순수 함수 3곳 수렴.
> fullend 전체 Toulmin 사용: 53
