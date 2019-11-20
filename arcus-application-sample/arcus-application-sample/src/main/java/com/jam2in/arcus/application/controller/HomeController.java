package com.jam2in.arcus.application.controller;

import com.jam2in.arcus.application.component.ArcusClientWrapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.*;

@Controller
public class HomeController {

    @Autowired
    private ArcusClientWrapper arcus;

    @RequestMapping("/")
    public String home(Model model) {

        String cacheKey = "kubernetes:arcus";
        Object cacheItem = arcus.get(cacheKey);

        if (cacheItem == null) {
            arcus.set(cacheKey, 10, "hello arcus!");
            cacheItem = "hello spring!";
        }

        model.addAttribute("cacheItem", cacheItem.toString());

        return "home";
    }

}
