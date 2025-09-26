#!/bin/sh
# ==== 系统信息 ====
ubus call system board                # 系统板卡信息 (型号/架构/版本)
ubus call system info                 # 系统运行状态 (uptime/load)
ubus call zwrt_bsp.thermal get_cpu_temp  # cpu 温度

# ==== 网络信息 ====
ubus call zte_nwinfo_api nwinfo_get_netinfo   # 获取 LTE/NR 信号、基站等信息
ubus call network.interface.lan status        # LAN 接口状态
ubus call network.interface.zte_wan status    # WAN IPv4 状态
ubus call network.interface.zte_wan6 status   # WAN IPv6 状态
ubus call network.wireless status             # 无线网络状态

# ==== Luci 相关 ====
# ubus call luci getUSBDevices          # USB 设备列表
# ubus call luci getBlockDevices        # 存储设备列表
# ubus call luci getProcessList         # 进程列表

# ==== 数据统计 ====
ubus call zwrt_data get_wwandst '{"source_module":"web","cid":1,"type":1,"clear_date_record":0}'  
# 实时/累计流量统计
ubus call zwrt_data get_wwandst_monthlimit '{"source_module":"web","cid":1}'  

# 月流量限制设置

# ==== 硬件状态 ====
ubus call zwrt_bsp.thermal list       # 温度传感器信息

# ==== 调试命令 ====
ubus list                             # 列出所有 service
ubus -v list <service>                # 查看某个 service 的详细方法 (比如 ubus -v list zte_nwinfo_api)

ubus call zwrt_bsp.usb set '{"mode":"debug"}'
ubus call zwrt_mc.device.manager get_device_info
ubus call zwrt_data get_wwandst '{"source_module": "web","cid": 1,"type": 4}'