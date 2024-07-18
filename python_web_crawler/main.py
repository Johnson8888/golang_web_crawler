# -*- coding: UTF-8 -*-
import json
import os
import time
import requests
import cv2
import json
import numpy as np
from PIL import Image
from bs4 import BeautifulSoup

from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions
from selenium.common.exceptions import TimeoutException
from selenium.webdriver.support.ui import WebDriverWait




class ProductMananger(object):

    def __init__(self):
        self.driver = webdriver.Chrome()
        super().__init__()
    
    def load_image_with_url_and_target_folder(self,url,folder):
        
        album_folder = folder + "/album"
        spec_folder = folder + "/spec"
        detail_folder = folder + "/detail"
        create_folders = [album_folder,spec_folder,detail_folder]
        for ff in create_folders:
            if os.path.exists(ff) == False:
                os.makedirs(ff)

        self.driver.get("http://www.1688.com")
        cookies_path = os.getcwd() + "/selenium/1688_cookies.txt"
        with open(cookies_path,'r') as f:
            cookies = json.loads(f.read())
            for cookie in cookies:
                break;
                # self.driver.add_cookie(cookie)
        self.driver.refresh()

        self.driver.get(url)

        time.sleep(20)

        iDetailConfig = self.driver.execute_script("return iDetailConfig;")
        detail_config_path = folder + "/config.json"
        with open(detail_config_path,'w') as f:
            f.write(json.dumps(iDetailConfig))
        print("iDetailConfig = {}".format(iDetailConfig))

        iDetailData = self.driver.execute_script("return iDetailData;")
        detail_data_path = folder + "/data.json"
        with open(detail_data_path,"w") as f:
            f.write(json.dumps(iDetailData))

        print("iDetailData = {}".format(iDetailData))


        if 'sku' in iDetailData and (isinstance(iDetailData['sku'],dict) == True) and ('skuProps' in iDetailData['sku']):
            
            sku_props = iDetailData['sku']['skuProps']
            sku_props_path = spec_folder + "/spec.json"
            with open(sku_props_path,"w") as f:
                f.write(json.dumps(sku_props))
            print("sku_props = {}".format(sku_props))

            if len(sku_props) > 0:
                props_value_list = sku_props[0]['value']
                for dic in props_value_list: 
                    if "imageUrl" in dic and "name" in dic:
                        imageUrl = dic['imageUrl'].replace("jpg","400x400.jpg")
                        name = dic['name'].replace("/","每")
                        local_file_path = spec_folder + "/" + name + ".jpg"
                        r = requests.get(imageUrl)
                        with open(local_file_path, 'wb') as f:
                            f.write(r.content)
        
        album_list = self.driver.find_elements_by_xpath("//div[@id='dt-tab']//li")
        for (index,album) in enumerate(album_list):
            album_data_imgs = album.get_attribute("data-imgs")
            if album_data_imgs == None:continue
            album_data_json = json.loads(album_data_imgs)
            album_url = album_data_json['preview']
            file_name = album_url.split("/")[-1]
            local_file_path = album_folder + "/" + str(index)  + "_" + file_name
            print (local_file_path)
            r = requests.get(album_url)
            with open(local_file_path, 'wb') as f:
                f.write(r.content)

            print("album_url = {} i = {}".format(album_data_json['preview'],index))
        
        # 查找详情的原理 按Ctrl+F调出查找，输入“加载中”三个字点下一个，会定位到一条代码上。所有代码中，有“加载中”的只有一条，所以不存在选错的情况。
        # //div[@id='desc-lazyload-container']
        desc_div = self.driver.find_element_by_xpath("//div[@id='desc-lazyload-container']")
        desc_url = desc_div.get_attribute("data-tfs-url")
        text = self.get_text(url=desc_url)
        self.run(text=text,target_folder=folder + "/detail")

        self.image_size(target_folder=folder)

        cookies = self.driver.get_cookies()
        save_file_path = os.getcwd() + "/1688_cookies.txt"
        cookie_string = json.dumps(cookies)
        with open(save_file_path,"w") as f:
            f.write(cookie_string)
        print(cookie_string)


    def get_detail_image_with_url(self,url,folder):
        '''
        获取详情图片
        '''
        text = self.get_text(url=url)
        self.run(text=text,target_folder=folder + "/detail")
        self.image_size(target_folder=folder)




    def get_scroll_image(self,url,target_folder):
        resp = requests.get(url)
        print(resp.text)

    def get_text(self,url):
        resp = requests.get(url)
        print(resp.text)
        return resp.text


    def run(self,text,target_folder):

        text = (text).replace("var offer_details=", "").replace(";", "")
        json_result = json.loads(text)

        obj = json_result['content']

        soup = BeautifulSoup(obj, 'html.parser')
        image_list = soup.find_all('img')

        for (index,image) in enumerate(image_list):
            image_url = image['src']
            file_name = image_url.split("/")[-1]
            print(file_name)
            local_file_path = target_folder + "/" + str(index)  + "_" + file_name
            print (local_file_path)
            r = requests.get(image_url)
            with open(local_file_path, 'wb') as f:
                f.write(r.content)

    def image_size(self,target_folder):
        '''
        修改图片的尺寸大小
        '''
        for parent, dirnames, filenames in os.walk(target_folder, followlinks=True):
            for filename in filenames:
                if filename.endswith("jpg") == False:
                    continue
                file_path = os.path.join(parent, filename)
                img = Image.open(file_path)
                # if img.size[0] > 400 or img.size[1] > 400:
                #     continue
                if img.size[0] == img.size[1] and img.size[0] < 480 and img.size[1] < 480:
                    image_source = cv2.imdecode(np.fromfile(file_path, dtype=np.uint8), cv2.IMREAD_UNCHANGED)
                    image = cv2.resize(image_source, (480, 480))
                    cv2.imencode('.jpg', image)[1].tofile(file_path)

                print (img.size)
                # print('文件名：%s' % filename)
                print('文件完整路径：%s\n' % file_path)



    def get_cookise(self):
        self.driver.get("http://www.1688.com")
        time.sleep(30)
        cookies = self.driver.get_cookies()
        print(cookies)
        save_file_path = os.getcwd() + "/1688_cookies.txt"
        cookie_string = json.dumps(cookies)
        with open(save_file_path,"w") as f:
            f.write(cookie_string)
        print(cookie_string)






if __name__ == '__main__':
    

    folder = "./594677944566"
    p = ProductMananger()
    url = "https://img.alicdn.com/tfscom/TB1Csv6p1bviK0jSZFNXXaApXXa"
    p.get_detail_image_with_url(
        url=url,
        folder=folder
    )


