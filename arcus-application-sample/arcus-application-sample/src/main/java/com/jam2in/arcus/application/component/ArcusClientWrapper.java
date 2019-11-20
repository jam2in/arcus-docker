package com.jam2in.arcus.application.component;

import net.spy.memcached.ArcusClient;
import net.spy.memcached.ConnectionFactoryBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

import javax.annotation.PostConstruct;
import javax.annotation.PreDestroy;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;

@Component
public class ArcusClientWrapper {

    public final static String ADDRESS;
    public final static String SERVICE_CODE;

    static {
        ADDRESS = System.getenv("ARCUS_ADDRESS");
        SERVICE_CODE = System.getenv("ARCUS_SERVICE_CODE");
    }

    private static final Logger log = LoggerFactory.getLogger(ArcusClientWrapper.class);

    private ArcusClient arcusClient;

    @PostConstruct
    public void postConstruct() {
        arcusClient = ArcusClient.createArcusClient(ADDRESS, SERVICE_CODE, new ConnectionFactoryBuilder());
    }

    @PreDestroy
    public void preDestroy() {
        if (arcusClient != null) {
            arcusClient.shutdown();
        }
    }

    @SuppressWarnings("unchecked")
    public <T> T get(String key) {
        Future<Object> future;
        Object result = null;

        try {
            future = arcusClient.asyncGet(key);
            log.debug("Get operation. (key=\"{}\")", key);
        } catch (Exception e) {
            log.error("error", e);
            return null;
        }

        try {
            result = future.get(700, TimeUnit.MILLISECONDS);
        } catch (Exception e) {
            log.error("error", e);
            future.cancel(true);
        }

        log.debug("Get operation result. (key=\"{}\", result={})", key, result != null);

        return (T) result;
    }

    public boolean set(String key, int expireTime, Object value) {

        Future<Boolean> future;
        boolean result = false;

        try {
            future = arcusClient.set(key, expireTime, value);
            log.debug("Set operation. (key=\"{}\", expireTime={})", key, expireTime);
        } catch (Exception e) {
            log.error("error", e);
            return false;
        }

        try {
            result = future.get(700, TimeUnit.MILLISECONDS);
        } catch (Exception e) {
            log.error("error", e);
            future.cancel(true);
        }

        log.debug("Set operation result. (key=\"{}\", result={})", key, result);

        return result;
    }


}
