-- 高危安全事件测试数据
-- 使用方法: mysql -u root -p fo_sentinel < test_high_severity_events.sql

USE fo_sentinel;

-- 设置时区为东八区
SET time_zone = '+08:00';

-- 1. Apache Struts 远程代码执行漏洞
INSERT INTO security_events (
    subscription_id, title, description, severity, status,
    source_url, event_time, cve_id, cvss_score,
    affected_vendor, affected_product, unique_hash, created_at, updated_at
) VALUES (
    1,
    'Apache Struts 2 远程代码执行漏洞 CVE-2023-50164',
    'Apache Struts 2.0.0 至 2.3.37、2.5.0 至 2.5.32 版本存在文件上传逻辑缺陷，攻击者可通过构造恶意请求实现远程代码执行，获取服务器控制权。该漏洞影响广泛，已被用于实际攻击。',
    'critical',
    'new',
    'https://nvd.nist.gov/vuln/detail/CVE-2023-50164',
    '2023-12-07 15:30:00',
    'CVE-2023-50164',
    9.8,
    'Apache',
    'Struts 2.0.0-2.5.32',
    MD5('CVE-2023-50164'),
    NOW(3),
    NOW(3)
);

-- 2. Spring Framework RCE 漏洞
INSERT INTO security_events (
    subscription_id, title, description, severity, status,
    source_url, event_time, cve_id, cvss_score,
    affected_vendor, affected_product, unique_hash, created_at, updated_at
) VALUES (
    1,
    'Spring Framework 远程代码执行漏洞 CVE-2022-22965 (Spring4Shell)',
    'Spring Framework 5.3.0-5.3.17、5.2.0-5.2.19 及更早版本在 JDK 9+ 环境下存在 RCE 漏洞。攻击者可通过特制 HTTP 请求修改 Tomcat 日志配置，写入 webshell 获取服务器权限。',
    'critical',
    'new',
    'https://nvd.nist.gov/vuln/detail/CVE-2022-22965',
    '2022-03-31 10:15:00',
    'CVE-2022-22965',
    9.8,
    'VMware/Spring',
    'Spring Framework 5.2.x-5.3.17',
    MD5('CVE-2022-22965'),
    NOW(3),
    NOW(3)
);

-- 3. Log4j 远程代码执行漏洞 (Log4Shell)
INSERT INTO security_events (
    subscription_id, title, description, severity, status,
    source_url, event_time, cve_id, cvss_score,
    affected_vendor, affected_product, unique_hash, created_at, updated_at
) VALUES (
    1,
    'Apache Log4j 远程代码执行漏洞 CVE-2021-44228 (Log4Shell)',
    'Apache Log4j 2.0-2.14.1 版本存在 JNDI 注入漏洞，攻击者可通过构造特殊的日志内容触发 LDAP/RMI 远程加载恶意类，实现任意代码执行。影响数百万 Java 应用，被评为近十年最严重漏洞之一。',
    'critical',
    'new',
    'https://nvd.nist.gov/vuln/detail/CVE-2021-44228',
    '2021-12-10 02:00:00',
    'CVE-2021-44228',
    10.0,
    'Apache',
    'Log4j 2.0-2.14.1',
    MD5('CVE-2021-44228'),
    NOW(3),
    NOW(3)
);

-- 4. Fastjson 反序列化漏洞
INSERT INTO security_events (
    subscription_id, title, description, severity, status,
    source_url, event_time, cve_id, cvss_score,
    affected_vendor, affected_product, unique_hash, created_at, updated_at
) VALUES (
    1,
    'Fastjson 反序列化远程代码执行漏洞 CVE-2022-25845',
    'Fastjson 1.2.80 及之前版本存在反序列化漏洞，攻击者可通过构造恶意 JSON 数据绕过 AutoType 安全机制，加载恶意类实现远程代码执行。该漏洞在国内 Java 应用中影响广泛。',
    'high',
    'new',
    'https://nvd.nist.gov/vuln/detail/CVE-2022-25845',
    '2022-02-23 18:20:00',
    'CVE-2022-25845',
    8.1,
    'Alibaba',
    'Fastjson <=1.2.80',
    MD5('CVE-2022-25845'),
    NOW(3),
    NOW(3)
);

-- 5. GitLab 认证绕过漏洞
INSERT INTO security_events (
    subscription_id, title, description, severity, status,
    source_url, event_time, cve_id, cvss_score,
    affected_vendor, affected_product, unique_hash, created_at, updated_at
) VALUES (
    1,
    'GitLab 身份认证绕过漏洞 CVE-2023-7028',
    'GitLab CE/EE 16.1.0-16.1.5、16.2.0-16.2.8、16.3.0-16.3.6 版本存在认证绕过漏洞，攻击者可通过发送密码重置邮件到未验证的邮箱地址，接管任意用户账户，包括管理员账户。',
    'critical',
    'new',
    'https://nvd.nist.gov/vuln/detail/CVE-2023-7028',
    '2024-01-11 09:45:00',
    'CVE-2023-7028',
    9.6,
    'GitLab',
    'GitLab CE/EE 16.1.0-16.3.6',
    MD5('CVE-2023-7028'),
    NOW(3),
    NOW(3)
);

