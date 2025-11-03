package com.turingq.eesdk;

/**
 * 证书主题
 */
public class Subject {
    private String commonName;
    private String[] country;
    private String[] organization;
    private String[] organizationalUnit;
    private String[] locality;
    private String[] province;
    
    public Subject(String commonName) {
        this.commonName = commonName;
    }
    
    public String getCommonName() { return commonName; }
    public void setCommonName(String commonName) { this.commonName = commonName; }
    
    public String[] getCountry() { return country; }
    public void setCountry(String[] country) { this.country = country; }
    
    public String[] getOrganization() { return organization; }
    public void setOrganization(String[] organization) { this.organization = organization; }
    
    public String[] getOrganizationalUnit() { return organizationalUnit; }
    public void setOrganizationalUnit(String[] organizationalUnit) { this.organizationalUnit = organizationalUnit; }
    
    public String[] getLocality() { return locality; }
    public void setLocality(String[] locality) { this.locality = locality; }
    
    public String[] getProvince() { return province; }
    public void setProvince(String[] province) { this.province = province; }
}
