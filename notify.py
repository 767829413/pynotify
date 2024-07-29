import pyautogui
import time
import cv2
import numpy as np
from datetime import datetime


def capture_tray_icon():
    # 获取屏幕分辨率
    screen_width, screen_height = pyautogui.size()

    # 定义要截取的区域（右下角）
    width, height = 313, 37
    x = screen_width - width
    y = screen_height - height
    region = (x, y, width, height)

    # 截取屏幕
    screenshot = pyautogui.screenshot(region=region)

    # 转换为OpenCV格式
    return cv2.cvtColor(np.array(screenshot), cv2.COLOR_RGB2BGR)


def detect_icon_change(prev_image, current_image, threshold=30):
    # 计算两图之间的差异
    diff = cv2.absdiff(prev_image, current_image)
    non_zero_count = np.count_nonzero(diff)

    # 如果差异像素数量超过阈值，则认为图标发生了变化
    return non_zero_count > threshold


def write_to_file(message):
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    with open("icon_changes.log", "a") as f:
        f.write(f"{timestamp}: {message}\n")


def main():
    prev_image = capture_tray_icon()
    while True:
        time.sleep(1)  # 每秒检查一次
        current_image = capture_tray_icon()

        if detect_icon_change(prev_image, current_image):
            write_to_file("企业微信图标闪动detected")
            print("检测到企业微信图标闪动")

        prev_image = current_image


if __name__ == "__main__":
    main()
