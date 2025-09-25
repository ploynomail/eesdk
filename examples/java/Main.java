import com.turingq.eesdk.EESdk;
import com.turingq.eesdk.Response;

/**
 * EE SDK Java使用示例
 */
public class Main {
    public static void main(String[] args) {
        System.out.println("=== EE SDK Java示例 ===");
        
        EESdk client = null;
        try {
            // 创建客户端 - 使用httpbin.org进行测试
            client = new EESdk(
                "http://httpbin.org",
                "your-api-key-here",
                30,  // 超时30秒
                3    // 重试3次
            );
            
            // 示例1: 健康检查
            System.out.println("\n=== 健康检查示例 ===");
            if (client.healthCheck()) {
                System.out.println("服务健康状态正常");
            } else {
                System.out.println("健康检查失败");
            }
            
            // 示例2: GET请求
            System.out.println("\n=== GET请求示例 ===");
            Response response = client.executeRequest("GET", "/get");
            printResponse(response);
            
            // 示例3: POST请求
            System.out.println("\n=== POST请求示例 ===");
            String jsonBody = "{\"name\":\"张三\",\"email\":\"zhangsan@example.com\",\"age\":25}";
            response = client.executeRequest("POST", "/post", jsonBody);
            printResponse(response);
            
            // 示例4: PUT请求
            System.out.println("\n=== PUT请求示例 ===");
            String updateBody = "{\"name\":\"李四\",\"age\":28}";
            response = client.executeRequest("PUT", "/put", updateBody);
            printResponse(response);
            
            // 示例5: 错误处理
            System.out.println("\n=== 错误处理示例 ===");
            response = client.executeRequest("GET", "/status/404");
            printResponse(response);
            
        } catch (Exception e) {
            System.err.println("SDK使用过程中发生错误: " + e.getMessage());
            e.printStackTrace();
        } finally {
            // 清理资源
            if (client != null) {
                client.destroy();
            }
        }
        
        System.out.println("\nJava示例执行完成");
    }
    
    private static void printResponse(Response response) {
        if (response.hasError()) {
            System.out.println("错误: " + response.getErrorMessage());
        } else {
            System.out.println("状态码: " + response.getStatusCode());
            System.out.println("响应体: " + response.getBody());
            System.out.println("请求成功: " + response.isSuccess());
        }
    }
}
