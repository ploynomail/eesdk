package com.turingq.eesdk;

/**
 * HTTP响应对象
 */
public class Response {
    private final int statusCode;
    private final String body;
    private final String errorMessage;
    
    public Response(int statusCode, String body, String errorMessage) {
        this.statusCode = statusCode;
        this.body = body;
        this.errorMessage = errorMessage;
    }
    
    /**
     * 获取状态码
     * 
     * @return HTTP状态码
     */
    public int getStatusCode() {
        return statusCode;
    }
    
    /**
     * 获取响应体
     * 
     * @return 响应体字符串
     */
    public String getBody() {
        return body;
    }
    
    /**
     * 获取错误信息
     * 
     * @return 错误信息，如果没有错误则为null
     */
    public String getErrorMessage() {
        return errorMessage;
    }
    
    /**
     * 检查是否有错误
     * 
     * @return 如果有错误返回true
     */
    public boolean hasError() {
        return errorMessage != null && !errorMessage.isEmpty();
    }
    
    /**
     * 检查是否成功（状态码在200-299范围内且无错误）
     * 
     * @return 如果成功返回true
     */
    public boolean isSuccess() {
        return !hasError() && statusCode >= 200 && statusCode < 300;
    }
    
    @Override
    public String toString() {
        return String.format("Response{statusCode=%d, body='%s', errorMessage='%s'}", 
                           statusCode, body, errorMessage);
    }
}
