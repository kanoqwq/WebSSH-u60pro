<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="title-section">
        <h1 class="title">
          <div class="uptime-value" style="font-size: 20px">
            运行中 {{ formatUptime(deviceInfo.device_uptime as any) }}
          </div>
        </h1>
        <div class="status-indicator">
          <div :class="['status-dot', dataReady ? 'online' : 'offline']"></div>
          <span class="status-text">{{ dataReady ? '已连接' : '未加载' }}</span>
        </div>
      </div>

      <div class="controls">
        <div style="display: flex; gap: 10px">
          <div style="display: flex; position: relative">
            <div class="signal-bars">
              <div
                v-for="n in 5"
                :key="n"
                :class="['bar full-bar', { active: n <= signalBars }]"></div>
            </div>
            <span
              style="position: absolute; font-size: 12px; top: -7px; left: -1px"
              >{{ networkType }}{{ is5GA ? 'A' : '' }}</span
            >
          </div>

          <div style="display: flex; align-items: center">
            <div
              :class="[
                'battery',
                {
                  charging:
                    deviceInfo?.bat_charger_connect &&
                    deviceInfo.bat_charger_connect == '1',
                },
              ]">
              <div
                class="battery-level"
                :style="{ width: `${deviceInfo.bat_percent || '0'}%` }"></div>
              <div class="battery-head"></div>
              <!-- 充电状态 -->
              <div class="battery-charging">⚡</div>
            </div>
            <span style="padding-left: 10px" v-if="deviceInfo.bat_percent"
              >{{ deviceInfo.bat_percent }} %</span
            >
          </div>
        </div>

        <div class="auto-refresh-controls">
          <button class="btn btn-primary" @click="refresh">刷新</button>
          <select
            v-model="refreshInterval"
            @change="updateRefreshInterval"
            :disabled="!autoRefresh">
            <option value="1000">1秒</option>
            <option value="2000">2秒</option>
            <option value="5000">5秒</option>
            <option value="10000">10秒</option>
          </select>
          <label class="checkbox-label">
            <input
              type="checkbox"
              v-model="autoRefresh"
              @change="toggleAutoRefresh" />
            <span class="checkmark"></span>
            自动
          </label>
        </div>
        <div class="auto-refresh-controls">
          <button
            class="btn btn-primary"
            @click="oneClickDebug"
            >一键ADB</button
          >
          <button
            class="btn btn-primary"
            @click="oneClickDebugClose"
            >关闭ADB</button
          >
          <button
            class="btn btn-primary"
            @click="smsForwardHandler"
            >短信转发</button
          >
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading && !dataReady" class="loading">
      <div class="loading-spinner"></div>
      <p>正在加载数据...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error">
      <div class="error-icon">⚠️</div>
      <h3>加载失败</h3>
      <p style="margin: 10px 0">{{ error }}</p>
      <button class="btn btn-danger" @click="refresh">重试</button>
    </div>

    <!-- 数据展示 -->
    <div v-else-if="dataReady" class="content">
      
      <!-- NR 5G 载波 -->
      <div class="card">
        <div class="card-header">
          <h3 class="hd">
            <img style="width: 24px" :src="NetworkIcon" alt="" />5G 载波信息
          </h3>
          <span v-if="networkType != '5G'" class="tag warning">未激活</span>
          <span v-else class="tag success">
            {{ networkType }}{{ is5GA ? 'A' : '' }}
            
            ({{ d.nr5g_action_band?.toUpperCase() ?? '-' }}{{
                formatNrca(d.nrca,'',0,3) != '-' ? ', N' + formatNrca(d.nrca,'',0,3) : '' }}{{
                formatNrca(d.nrca,'',1,3) != '-' ? ', N' + formatNrca(d.nrca,'',1,3) : '' }})
            </span>
        </div>
        <div class="card-content">
          <div class="signal-grid">
            <table class="mytable" width="100%">
              <tr>
                <td width="13%"></td>
                <td width="9%">PCI</td>
                <td width="11%">5G<br/>频段</td>
                <td width="16%">5G<br/>频点</td>
                <td width="11%">DL<br/>带宽</td>
                <td width="10%">RSRP</td>
                <td width="10%">RSRQ</td>
                <td width="10%">SINR</td>
                <td width="10%">RSSI</td>
              </tr>
              <tr>
                <td>PCC</td>
                <td>{{ d.nr5g_pci ?? '-' }}</td>
                <td>{{ d.nr5g_action_band?.toUpperCase() ?? '-' }}</td>
                <td>{{ d.nr5g_action_channel ?? '-' }}</td>
                <td>{{ d.nr5g_bandwidth ?? '-' }}</td>
                <td class="dbmstyle">{{ d.nr5g_rsrp }}</td>
                <td>{{ d.nr5g_rsrq }}</td>
                <td>{{ d.nr5g_snr }}</td>
                <td class="dbmstyle">{{ d.nr5g_rssi }}</td>
              </tr>
              <tr>
                <td>SCC0</td>
                <td>{{ formatNrca(d.nrca,'',0,1) }}</td>
                <td>{{ formatNrca(d.nrca,'N',0,3) }}</td>
                <td>{{ formatNrca(d.nrca,'',0,4) }}</td>
                <td>{{ formatNrca(d.nrca,'',0,5) }}</td>
                <td class="dbmstyle">{{ formatNrca(d.nrca,'',0,7) }}</td>
                <td>{{ formatNrca(d.nrca,'',0,8) }}</td>
                <td>{{ formatNrca(d.nrca,'',0,9) }}</td>
                <td class="dbmstyle">{{ formatNrca(d.nrca,'',0,10) }}</td>
              </tr>
              <tr>
                <td>SCC1</td>
                <td>{{ formatNrca(d.nrca,'',1,1) }}</td>
                <td>{{ formatNrca(d.nrca,'N',1,3) }}</td>
                <td>{{ formatNrca(d.nrca,'',1,4) }}</td>
                <td>{{ formatNrca(d.nrca,'',1,5) }}</td>
                <td class="dbmstyle">{{ formatNrca(d.nrca,'',1,7) }}</td>
                <td>{{ formatNrca(d.nrca,'',1,8) }}</td>
                <td>{{ formatNrca(d.nrca,'',1,9) }}</td>
                <td class="dbmstyle">{{ formatNrca(d.nrca,'',1,10) }}</td>
              </tr>
            </table>
          </div>
        </div>
      </div>

      <!-- NR 4G 载波 -->
      <div class="card">
        <div class="card-header">
          <h3 class="hd">
            <img style="width: 24px" :src="NetworkIcon" alt="" />4G 载波信息
          </h3>
          <span v-if="networkType != '4G'" class="tag warning">未激活</span>
          <span v-else class="tag success">
            {{ networkType }}{{ is5GA ? '+' : '' }}
            ({{ formatNrca(d.lteca,'B',0,1) }}{{
                formatNrca(d.lteca,'',1,1) != '-' ? ', B' + formatNrca(d.lteca,'',1,1) : '' }}{{
                formatNrca(d.lteca,'',2,1) != '-' ? ', B' + formatNrca(d.lteca,'',2,1) : '' }}{{
                formatNrca(d.lteca,'',3,1) != '-' ? ', B' + formatNrca(d.lteca,'',3,1) : '' }})
          </span>
        </div>
        <div class="card-content">
          <div class="signal-grid">
            <table class="mytable" width="100%">
              <tr>
                <td width="13%"></td>
                <td width="9%">PCI</td>
                <td width="11%">4G<br/>频段</td>
                <td width="16%">4G<br/>信道</td>
                <td width="11%">DL<br/>带宽</td>
                <td width="10%">RSRP</td>
                <td width="10%">RSRQ</td>
                <td width="10%">SINR</td>
                <td width="10%">RSSI</td>
              </tr>
              <tr>
                <td>PCC</td>
                <td>{{ d.lte_pci ?? '-' }}</td>
                <td>{{ formatNrca(d.lteca,'B',0,1) }}</td>
                <td>{{ d.wan_active_channel ?? '-' }}</td>
                <td>{{ formatNrca(d.lteca,'',0,4) }}</td>
                <td class="dbmstyle">{{ d.lte_rsrp }}</td>
                <td>{{ d.lte_rsrq }}</td>
                <td>{{ d.lte_snr }}</td>
                <td class="dbmstyle">{{ d.lte_rssi }}</td>
              </tr>
              <tr>
                <td>SCC0</td>
                <td>{{ formatNrca(d.lteca,'',1,0) }}</td>
                <td>{{ formatNrca(d.lteca,'B',1,1) }}</td>
                <td>{{ formatNrca(d.lteca,'',1,3) }}</td>
                <td>{{ formatNrca(d.lteca,'',1,4) }}</td>
                <td class="dbmstyle">{{ formatNrca(d.ltecasig,'',0,0) }}</td>
                <td>{{ formatNrca(d.ltecasig,'',0,1) }}</td>
                <td>{{ formatNrca(d.ltecasig,'',0,2) }}</td>
                <td class="dbmstyle">{{ formatNrca(d.ltecasig,'',0,3) }}</td>
              </tr>
              <tr>
                <td>SCC1</td>
                <td>{{ formatNrca(d.lteca,'',2,0) }}</td>
                <td>{{ formatNrca(d.lteca,'B',2,1) }}</td>
                <td>{{ formatNrca(d.lteca,'',2,3) }}</td>
                <td>{{ formatNrca(d.lteca,'',2,4) }}</td>
                <td class="dbmstyle">{{ formatNrca(d.ltecasig,'',1,0) }}</td>
                <td>{{ formatNrca(d.ltecasig,'',1,1) }}</td>
                <td>{{ formatNrca(d.ltecasig,'',1,2) }}</td>
                <td class="dbmstyle">{{ formatNrca(d.ltecasig,'',1,3) }}</td>
              </tr>
              <tr>
                <td>SCC2</td>
                <td>{{ formatNrca(d.lteca,'',3,0) }}</td>
                <td>{{ formatNrca(d.lteca,'B',3,1) }}</td>
                <td>{{ formatNrca(d.lteca,'',3,3) }}</td>
                <td>{{ formatNrca(d.lteca,'',3,4) }}</td>
                <td class="dbmstyle">{{ formatNrca(d.ltecasig,'',2,0) }}</td>
                <td>{{ formatNrca(d.ltecasig,'',2,1) }}</td>
                <td>{{ formatNrca(d.ltecasig,'',2,2) }}</td>
                <td class="dbmstyle">{{ formatNrca(d.ltecasig,'',2,3) }}</td>
              </tr>
            </table>
          </div>
        </div>
      </div>
      
      <!-- NR 5G信号卡片 -->
      <div class="card">
        <div class="card-header">
          <h3 class="hd">
            <img style="width: 24px" :src="NetworkIcon" alt="" />NR 5G 信号
          </h3>
          <span v-if="networkType != '5G'" class="tag warning">未激活</span>
          <span v-else class="tag success">已激活</span>
        </div>
        <div class="card-content">
          <div class="signal-grid">
            <div class="signal-item">
              <span class="label">RSRP</span>
              <div class="progress-bar">
                <div
                  class="progress-fill"
                  :style="{ width: getRsrpPercent(d.nr5g_rsrp) + '%' }"></div>
                <span class="progress-text">{{ formatDbm(d.nr5g_rsrp) }}</span>
              </div>
            </div>
            <div class="signal-item">
              <span class="label">RSRQ</span>
              <div class="progress-bar">
                <div
                  class="progress-fill"
                  :style="{ width: getRsrqPercent(d.nr5g_rsrq) + '%' }"></div>
                <span class="progress-text">{{ formatDb(d.nr5g_rsrq) }}</span>
              </div>
            </div>
            <div class="signal-item">
              <span class="label">SINR</span>
              <div class="progress-bar">
                <div
                  class="progress-fill"
                  :style="{ width: getSnrPercent(d.nr5g_snr) + '%' }"></div>
                <span class="progress-text">{{ formatSnr(d.nr5g_snr) }}</span>
              </div>
            </div>
            <div class="signal-item">
              <span class="label">RSSI</span>
              <div class="progress-bar">
                <div
                  class="progress-fill"
                  :style="{ width: getRssiPercent(d.nr5g_rssi) + '%' }"></div>
                <span class="progress-text">{{ formatDbm(d.nr5g_rssi) }}</span>
              </div>
            </div>
            <div class="signal-item">
              <span class="label">PCI</span>
              <span class="value">{{ d.nr5g_pci ?? '-' }}</span>
            </div>
            
            <div class="signal-item" width="100%">
              <span class="label">Cell ID</span>
              <span class="value">{{ d.nr5g_cell_id ?? '-' }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- LTE信号卡片 -->
      <div class="card">
        <div class="card-header">
          <h3 class="hd">
            <img style="width: 24px" :src="NetworkIcon" alt="" />LTE 信号
          </h3>
          <span v-if="networkType != '4G'" class="tag warning">未激活</span>
          <span v-else class="tag success">已激活</span>
        </div>
        <div class="card-content">
          <div class="signal-grid">
            
            <div class="signal-item">
              <span class="label">RSRP</span>
              <div class="progress-bar">
                <div
                  class="progress-fill"
                  :style="{ width: getRsrpPercent(d.lte_rsrp) + '%' }"></div>
                <span class="progress-text">{{ formatDbm(d.lte_rsrp) }}</span>
              </div>
            </div>
            <div class="signal-item">
              <span class="label">RSRQ</span>
              <div class="progress-bar">
                <div
                  class="progress-fill"
                  :style="{ width: getRsrqPercent(d.lte_rsrq) + '%' }"></div>
                <span class="progress-text">{{ formatDb(d.lte_rsrq) }}</span>
              </div>
            </div>
            <div class="signal-item">
              <span class="label">SINR</span>
              <div class="progress-bar">
                <div
                  class="progress-fill"
                  :style="{ width: getSnrPercent(d.lte_snr) + '%' }"></div>
                <span class="progress-text">{{ formatSnr(d.lte_snr) }}</span>
              </div>
            </div>
            <div class="signal-item">
              <span class="label">RSSI</span>
              <div class="progress-bar">
                <div
                  class="progress-fill"
                  :style="{ width: getRssiPercent(d.lte_rssi) + '%' }"></div>
                <span class="progress-text">{{ formatDbm(d.lte_rssi) }}</span>
              </div>
            </div>
            <div class="signal-item">
              <span class="label">PCI</span>
              <span class="value">{{ d.lte_pci ?? '-' }}</span>
            </div>
            <div class="signal-item">
              <span class="label">Cell ID</span>
              <span class="value">{{ d.cell_id ?? '-' }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 设备信息卡片 -->
      <div class="card device-info-card">
        <div class="card-header">
          <h3 class="hd">
            <img style="width: 24px" :src="DashboardIcon" alt="" />设备信息
          </h3>
        </div>
        <div class="card-content">
          <div class="device-stats">
            <div class="device-item">
              <div class="device-label">CPU 温度</div>
              <div class="temp-info">
                <div class="temp-value">
                  {{ formatCpuTemp(cpuTemp.cpuss_temp) }}
                </div>
                <div class="temp-bar">
                  <div
                    class="temp-fill"
                    :style="{
                      width: cpuTemp.cpuss_temp + '%',
                    }"></div>
                </div>
              </div>
            </div>

            <div class="device-item">
              <div class="device-label">
                电池温度<span
                  v-if="!(deviceInfo?.hightemp_datalimit_status == '0')"
                  style="color: orange"
                  >（温度保护）</span
                >
              </div>
              <div class="temp-info">
                <div class="temp-value">
                  {{ formatCpuTemp(deviceInfo.bat_temperature as any) }}
                </div>
                <div class="temp-bar">
                  <div
                    class="temp-fill"
                    :style="{
                      width: deviceInfo.bat_temperature + '%',
                    }"></div>
                </div>
              </div>
            </div>

            <div class="device-item">
              <div class="device-label">内存使用</div>
              <div class="memory-info">
                <div class="memory-details">
                  <span class="memory-used">{{
                    formatMemory(
                      ((deviceInfo.meminfo?.total || 0) as any) -
                        ((deviceInfo.meminfo?.free as any) || 0)
                    )
                  }}</span>
                  <span class="memory-separator">/</span>
                  <span class="memory-total">{{
                    formatMemory((deviceInfo.meminfo?.total as any) || 0)
                  }}</span>
                </div>
                <div class="memory-bar">
                  <div
                    class="memory-fill"
                    :style="{
                      width:
                        formatMemoryPercent(
                          (deviceInfo.meminfo?.total as any || 0) -
                            (deviceInfo.meminfo?.free as any || 0),
                          deviceInfo.meminfo?.total as any || 1
                        ) + '%',
                    }"></div>
                  <span class="memory-text"
                    >{{
                      formatMemoryPercent(
                        ((deviceInfo.meminfo?.total as any) || 0) -
                          ((deviceInfo.meminfo?.free as any) || 0),
                        (deviceInfo.meminfo?.total as any) || 1
                      )
                    }}%</span
                  >
                </div>
              </div>
            </div>

            <div class="device-item">
              <div class="device-label">
                CPU 负载
                {{
                  (100 - (deviceInfo.cpuinfo?.[0]?.idle as any) || 0).toFixed(2)
                }}
                %
              </div>
              <div class="device-values">
                <div class="load-item">
                  <span class="load-label">核心1</span>
                  <span class="load-value"
                    >{{
                      (
                        100 - (deviceInfo.cpuinfo?.[1]?.idle as any) || 0
                      ).toFixed(0)
                    }}
                    %</span
                  >
                </div>
                <div class="load-item">
                  <span class="load-label">核心2</span>
                  <span class="load-value"
                    >{{
                      (
                        100 - (deviceInfo.cpuinfo?.[2]?.idle as any) || 0
                      ).toFixed(0)
                    }}
                    %</span
                  >
                </div>
                <div class="load-item">
                  <span class="load-label">核心3</span>
                  <span class="load-value"
                    >{{
                      (
                        100 - (deviceInfo.cpuinfo?.[3]?.idle as any) || 0
                      ).toFixed(0)
                    }}
                    %</span
                  >
                </div>
                <div class="load-item">
                  <span class="load-label">核心4</span>
                  <span class="load-value"
                    >{{
                      (
                        100 - (deviceInfo.cpuinfo?.[4]?.idle as any) || 0
                      ).toFixed(0)
                    }}
                    %</span
                  >
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 网络信息卡片 -->
      <div class="card">
        <div class="card-header">
          <h3 class="hd">
            <img style="width: 24px" :src="InternetIcon" alt="" />网络信息
          </h3>
        </div>
        <div class="card-content">
          <div class="info-grid">
            <div class="info-item">
              <span class="label">运营商</span>
              <span class="value">{{
                d.network_provider_fullname || d.network_provider || '-'
              }}</span>
            </div>
            <div class="info-item">
              <span class="label">网络类型</span>
              <span class="value">{{ d.network_type || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">驻网状态</span>
              <span class="value">{{ d.simcard_roam || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">选择模式</span>
              <span class="value">{{ d.net_select_mode || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">选择策略</span>
              <span class="value">{{ d.net_select || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">信号强度</span>
              <div class="signal-bars">
                <div
                  v-for="n in 5"
                  :key="n"
                  :class="['bar', { active: n <= signalBars }]"></div>
              </div>
            </div>
            <div class="info-item">
              <span class="label">连接数量</span>
              <span class="value"
                >有线：{{ lanUserList?.lan_num || '-' }} / 无线：{{
                  lanUserList?.wireless_num || '-'
                }}</span
              >
            </div>
          </div>
        </div>
      </div>

      <!-- 流量统计卡片 -->
      <div class="card">
        <div class="card-header">
          <h3 class="hd">
            <img style="width: 24px" :src="ChartIcon" alt="" />流量统计
          </h3>
        </div>
        <div class="card-content">
          <div class="traffic-stats">
            <div class="traffic-item">
              <div class="traffic-label">上传速度</div>
              <div class="traffic-value upload">
                {{ formatSpeed(trafficData.real_tx_speed) }}
              </div>
            </div>
            <div class="traffic-item">
              <div class="traffic-label">下载速度</div>
              <div class="traffic-value download">
                {{ formatSpeed(trafficData.real_rx_speed) }}
              </div>
            </div>
            <div class="traffic-item">
              <div class="traffic-label">上传用量</div>
              <div class="traffic-value" style="font-size: 12px">
                <!-- {{ formatBytes(trafficData.real_tx_bytes) }} -->
                <div>当日：{{ formatBytes(trafficData.day_tx_bytes) }}</div>
                <div>本月：{{ formatBytes(trafficData.month_tx_bytes) }}</div>
              </div>
            </div>
            <div class="traffic-item">
              <div class="traffic-label">下载用量</div>
              <div class="traffic-value"></div>
              <div class="traffic-value" style="font-size: 12px">
                <!-- {{ formatBytes(trafficData.real_rx_bytes) }} -->
                <div>当日：{{ formatBytes(trafficData.day_rx_bytes) }}</div>
                <div>本月：{{ formatBytes(trafficData.month_rx_bytes) }}</div>
              </div>
            </div>
            <div class="traffic-item">
              <div class="traffic-label">总上传</div>
              <div class="traffic-value">
                {{ formatBytes(trafficData.total_tx_bytes as any) }}
              </div>
            </div>
            <div class="traffic-item">
              <div class="traffic-label">总下载</div>
              <div class="traffic-value">
                {{ formatBytes(trafficData.total_rx_bytes as any) }}
              </div>
            </div>
            <div class="traffic-item">
              <div class="traffic-label">最大上传速度</div>
              <div class="traffic-value">
                {{ formatSpeed(trafficData.real_max_tx_speed) }}
              </div>
            </div>
            <div class="traffic-item">
              <div class="traffic-label">最大下载速度</div>
              <div class="traffic-value">
                {{ formatSpeed(trafficData.real_max_rx_speed) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 频段与锁定卡片 -->
      <div class="card">
        <div class="card-header">
          <h3 class="hd">
            <img style="width: 24px" :src="LockIcon" alt="" />频段与锁定
          </h3>
        </div>
        <div class="card-content">
          <div class="info-grid">
            <div class="info-item">
              <span class="label">主载波</span>
              <span class="value">{{
                d.wan_active_band?.toUpperCase() || '-'
              }}</span>
            </div>
            <div class="info-item">
              <span class="label">工作频段</span>
              <span class="value"
                >{{
                  d.wan_active_band?.toUpperCase()
                    ? d.wan_active_band.toUpperCase() + ', '
                    : ''
                }}{{ currentActiveBands || '-' }}</span
              >
            </div>
            <div class="info-item">
              <span class="label">频道</span>
              <span class="value">{{ d.nr5g_action_channel ?? '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">带宽</span>
              <span class="value">{{
                d.nr5g_bandwidth ? d.nr5g_bandwidth + ' MHz' : '-'
              }}</span>
            </div>
            <div class="info-item">
              <span class="label">LTE 锁频</span>
              <span class="value">{{ d.lte_band_lock || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">NR SA 锁频</span>
              <span class="value">{{ d.nr5g_sa_band_lock || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">LTE 频段</span>
              <span
                class="value"
                style="
                  white-space: pre-wrap;
                  word-wrap: break-word;
                  overflow: hidden;
                "
                >{{ d.lte_band || '-' }}</span
              >
            </div>
          </div>
        </div>
      </div>

      <!-- 网络接口状态卡片 -->
      <div class="card">
        <div class="card-header">
          <h3 class="hd">
            <img style="width: 24px" :src="InterfaceIcon" alt="" />网络接口状态
          </h3>
        </div>
        <div class="card-content">
          <div class="interface-grid">
            <div class="interface-section">
              <h4>LAN 接口</h4>
              <div class="info-grid-compact">
                <div class="info-item">
                  <span class="label">IP 地址</span>
                  <span class="value">{{
                    lanData.ipv4_address?.[0]?.address || '-'
                  }}</span>
                </div>
                <div class="info-item">
                  <span class="label">网关</span>
                  <span class="value">{{
                    lanData.route?.[0]?.nexthop || '-'
                  }}</span>
                </div>
                <div class="info-item">
                  <span class="label">DNS 服务器</span>
                  <span class="value">{{
                    lanData['dns-server']?.join(', ') || '-'
                  }}</span>
                </div>
                <div class="info-item">
                  <span class="label">运行时间</span>
                  <span class="value">{{ formatUptime(lanData.uptime) }}</span>
                </div>
              </div>
            </div>

            <div class="interface-section">
              <h4>WAN IPv4 接口</h4>
              <div class="info-grid-compact">
                <div class="info-item">
                  <span class="label">IP 地址</span>
                  <span class="value">{{ wwanInfo?.ipv4_address || '-' }}</span>
                </div>
                <div class="info-item">
                  <span class="label">网关</span>
                  <span class="value">{{ wwanInfo?.ipv4_gateway || '-' }}</span>
                </div>
                <div class="info-item">
                  <span class="label">DNS 服务器</span>
                  <span class="value">{{
                    wanData['dns-server']?.join(', ') || '-'
                  }}</span>
                </div>
                <div class="info-item">
                  <span class="label">运行时间</span>
                  <span class="value">{{ formatUptime(wanData.uptime) }}</span>
                </div>
              </div>
            </div>

            <div class="interface-section">
              <h4>WAN IPv6 接口</h4>
              <div class="info-grid-compact">
                <div class="info-item">
                  <span class="label">IPv6 地址</span>
                  <span class="value">{{ wwanInfo?.ipv6_address || '-' }}</span>
                </div>
                <div class="info-item">
                  <span class="label">网关</span>
                  <span class="value">{{ wwanInfo?.ipv6_gateway || '-' }}</span>
                </div>
                <div class="info-item">
                  <span class="label">DNS 服务器</span>
                  <span class="value">{{
                    wan6Data['dns-server']?.join(', ') || '-'
                  }}</span>
                </div>
                <div class="info-item">
                  <span class="label">运行时间</span>
                  <span class="value">{{ formatUptime(wan6Data.uptime) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 标识信息卡片 -->
      <div class="card">
        <div class="card-header">
          <h3 class="hd">
            <img style="width: 24px" :src="TagIcon" alt="" />标识信息
          </h3>
        </div>
        <div class="card-content">
          <div class="info-grid">
            <div class="info-item">
              <span class="label">ICCID</span>
              <span class="value">{{ simInfo2.sim_iccid ?? '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">IMSI</span>
              <span class="value">{{ simInfo2.sim_imsi ?? '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">IMEI</span>
              <span class="value">{{ simInfo?.values?.imei ?? '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">Lock Status</span>
              <span class="value">{{
                simInfo?.values?.lock_status ?? '-'
              }}</span>
            </div>
            <div class="info-item">
              <span class="label">Modem MSN</span>
              <span class="value">{{ simInfo?.values?.modem_msn ?? '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">WLAN MAC</span>
              <span class="value">{{
                simInfo?.values?.wlan_mac_address ?? '-'
              }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty">
      <div class="empty-icon">📱</div>
      <h3>暂无数据</h3>
      <p>请点击刷新按钮获取设备信息</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import ChartIcon from '@/assets/svgs/chart.svg';
import DashboardIcon from '@/assets/svgs/dashboard.svg';
import InterfaceIcon from '@/assets/svgs/interface.svg';
import InternetIcon from '@/assets/svgs/internet.svg';
import LockIcon from '@/assets/svgs/lock.svg';
import NetworkIcon from '@/assets/svgs/network.svg';
import TagIcon from '@/assets/svgs/tag.svg';

import axios from 'axios';
import { ElMessage, ElNotification } from 'element-plus';
import { computed, onMounted, onUnmounted, ref } from 'vue';

interface UbusResponse<T = any> {
  code: number;
  msg: string;
  result?: T;
}

interface NetInfoResult {
  [key: string]: any;
}

interface NetworkInterface {
  up: boolean;
  device?: string;
  proto?: string;
  uptime?: number;
  ipv4_address?: Array<{ address: string; mask: number }>;
  ipv6_address?: Array<{ address: string; mask: number }>;
  route?: Array<{ nexthop: string }>;
  'dns-server'?: string[];
}

interface TrafficData {
  real_tx_speed: number;
  real_rx_speed: number;
  real_tx_bytes: number;
  real_rx_bytes: number;
  real_max_tx_speed: number;
  real_max_rx_speed: number;
  total_tx_bytes: number;
  total_rx_bytes: number;
  day_tx_bytes: number;
  day_rx_bytes: number;
  month_tx_bytes: number;
  month_rx_bytes: number;
}

interface SystemInfo {
  localtime: number;
  uptime: number;
  load: number[];
  memory: {
    total: number;
    free: number;
    shared: number;
    buffered: number;
    available: number;
    cached: number;
  };
  root: {
    total: number;
    free: number;
    used: number;
    avail: number;
  };
  tmp: {
    total: number;
    free: number;
    used: number;
    avail: number;
  };
  swap: {
    total: number;
    free: number;
  };
}

interface DeviceInfo {
  hightemp_datalimit_status: string;
  quicken_power_on: string;
  bat_online: string;
  bat_health: string;
  bat_mode: string;
  bat_low_power: string;
  bat_percent: string;
  bat_level: string;
  bat_temperature: string;
  bat_charger_connect: string;
  bat_charger_type: string;
  bat_charger_status: string;
  bat_ui_charger_type: string;
  bat_temperature_level: string;
  external_charging_flag: string;
  bat_time_to_full: string;
  bat_time_to_empty: string;
  power_adapter: string;
  device_uptime: string;
  cpuinfo: {
    name: string;
    idle: string;
    gnice?: string;
  }[];
  meminfo: {
    total: string;
    free: string;
    avaliable: string;
  };
  flashinfo: {
    filesystem: string;
    size: string;
    used: string;
    avail: string;
    use: string;
    mounted_on: string;
  }[];
}

interface CpuTemp {
  cpuss_temp: number;
}

interface SimInfo {
  values: {
    digitalcode: string;
    imei: string;
    imei2: string;
    lock_status: string;
    modem_msn: string;
    wlan_mac_address: string;
  };
}

interface SimInfo2 {
  sim_iccid: string;
  sim_imsi: string;
}

interface WwanInfo {
  connect_fail_count: 0;
  connect_status: string;
  ipv4_address: string;
  ipv4_dev_name: string;
  ipv4_dns_prefer: string;
  ipv4_dns_standby: string;
  ipv4_gateway: string;
  ipv4_netmask: string;
  ipv6_address: string;
  ipv6_dev_name: string;
  ipv6_dns_prefer: string;
  ipv6_dns_standby: string;
  ipv6_gateway: string;
  roam_enable: number;
}

interface LanUserList {
  access_total_num: number;
  lan_num: number;
  wireless_num: number;
  offline_num: number;
  guest_num_24g: number;
  guest_num_5g: number;
  guest_num_6g: number;
}

// 响应式数据
const loading = ref(false);
const error = ref<string | null>(null);
const data = ref<NetInfoResult | null>(null);
const lanData = ref<NetworkInterface>({} as NetworkInterface);
const wanData = ref<NetworkInterface>({} as NetworkInterface);
const wan6Data = ref<NetworkInterface>({} as NetworkInterface);
const trafficData = ref<TrafficData>({} as TrafficData);
const deviceInfo = ref<DeviceInfo>({} as DeviceInfo);
const cpuTemp = ref<CpuTemp>({} as CpuTemp);
const simInfo = ref<SimInfo>({} as SimInfo);
const simInfo2 = ref<SimInfo2>({} as SimInfo2);
const wwanInfo = ref<WwanInfo>({} as WwanInfo);
const lanUserList = ref<LanUserList>({} as LanUserList);
const networkType = computed(() => {
  const val = d.value?.network_type;
  if (
    val?.includes('NR') ||
    val?.includes('5G') ||
    val?.includes('SA') ||
    val?.includes('NSA') ||
    val?.includes('ENDC')
  )
    return '5G';
  if (val?.includes('4G') || val?.includes('LTE')) return '4G';
  if (val?.includes('HSPA')) return 'H+';
  if (val?.includes('3G')) return '3G';
  return '';
});
const currentActiveBands = computed(() => {
  const val = networkType.value != '5G' ? data.value?.lteca : data.value?.nrca;
  if (!val) return null;
  const list = (val as string).split(';').filter((el) => el);
  if (list.length == 0) return null;
  if (networkType.value != '5G' && list.length == 1) return null;

  const res = list
    .map((item) =>
      networkType.value != '5G'
        ? item[1]
          ? item.split(',')[1]
          : ''
        : item[3]
        ? 'N' + item.split(',')[3]
        : ''
    )
    .join(', ')
    .replace(/,/g, ',');
  return res;
});

const is5GA = computed(() => {
  if (!currentActiveBands.value) return false;
  return currentActiveBands.value.split(',').length >= 2;
});

// 自动刷新控制
const autoRefresh = ref(true);
const refreshInterval = ref(1000);
let refreshTimer: number | null = null;

// 请求体定义
const netInfoRequest = {
  service: 'zte_nwinfo_api',
  method: 'nwinfo_get_netinfo',
  params: {},
};

const lanRequest = {
  service: 'network.interface.lan',
  method: 'status',
  params: {},
};

const wanRequest = {
  service: 'network.interface.zte_wan',
  method: 'status',
  params: {},
};

const wan6Request = {
  service: 'network.interface.zte_wan6',
  method: 'status',
  params: {},
};

const trafficRequest = {
  service: 'zwrt_data',
  method: 'get_wwandst',
  params: { source_module: 'web', cid: 1, type: 4 },
};

const simInfoRequest = {
  service: 'uci',
  method: 'get',
  params: {
    config: 'zwrt_zte_mdm',
    section: 'device_info',
  },
};

const simInfo2Request = {
  service: 'zwrt_zte_mdm.api',
  method: 'get_sim_info',
  params: {},
};

const deviceInfoRequest = {
  service: 'zwrt_mc.device.manager',
  method: 'get_device_info',
  params: {},
};

const cpuTempRequest = {
  service: 'zwrt_bsp.thermal',
  method: 'get_cpu_temp',
  params: {},
};

const wwanRequest = {
  service: 'zwrt_data',
  method: 'get_wwaniface',
  params: {
    source_module: 'web',
    cid: 1,
    connect_status: '',
  },
};

const lanUserListRequest = {
  service: 'zwrt_router.api',
  method: 'router_get_user_list_num',
  params: {},
};

const openAdbRequest = {
  service: 'zwrt_bsp.usb',
  method: 'set',
  params: {
    mode: 'debug',
  },
};

const closeAdbRequest = {
  service: 'zwrt_bsp.usb',
  method: 'set',
  params: {
    mode: 'user',
  },
};

// 计算属性
const dataReady = computed(() => !!data.value);
const d = computed(() => data.value || {});

const signalBars = computed(() => {
  const bars = Number(d.value.signalbar || 0);
  if (Number.isNaN(bars)) return 0;
  return Math.max(0, Math.min(5, bars));
});

// 格式化函数
function formatDbm(v: unknown): string {
  const n = Number(v);
  if (Number.isNaN(n)) return '-';
  return `${n} dBm`;
}

function formatDb(v: unknown): string {
  const n = Number(v);
  if (Number.isNaN(n)) return '-';
  return `${n} dB`;
}

function formatSnr(v: unknown): string {
  const n = Number(v);
  if (Number.isNaN(n)) return '-';
  return `${n.toFixed(1)} dB`;
}

function formatUptime(seconds?: number): string {
  if (!seconds) return '-';
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;
  return `${hours}:${minutes}:${secs}`;
}

function formatSpeed(bytesPerSecond: number): string {
  if (!bytesPerSecond) return '0 B/s';
  const units = ['B/s', 'KB/s', 'MB/s', 'GB/s'];
  let size = bytesPerSecond;
  let unitIndex = 0;

  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }

  return `${size.toFixed(unitIndex === 0 ? 0 : 1)} ${units[unitIndex]}`;
}

function formatBytes(bytes: number): string {
  if (!bytes) return '-';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let size = bytes;
  let unitIndex = 0;

  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }

  return `${size.toFixed(unitIndex === 0 ? 0 : 1)} ${units[unitIndex]}`;
}

function formatNumber(num: number): string {
  if (!num) return '-';
  return num.toLocaleString();
}

// 系统信息格式化函数
function formatLoad(load: number): string {
  return (load / 1000).toFixed(2);
}

function formatMemory(KB: number): string {
  if (!KB) return '-';
  const units = ['B', 'KB', 'MB', 'GB'];
  let size = KB * 1024;
  let unitIndex = 0;

  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }

  return `${size.toFixed(unitIndex === 0 ? 0 : 1)} ${units[unitIndex]}`;
}

function formatMemoryPercent(used: number, total: number): number {
  if (!total) return 0;
  return Math.round((used / total) * 100);
}

function formatCpuTemp(temp: number): string {
  if (!temp) return '-';
  return `${temp}°C`;
}

// 信号强度百分比计算函数
function getRssiPercent(rssi: number): number {
  if (!rssi) return 0;
  // RSSI: -120dBm (0%) 到 -30dBm (100%)
  const min = -120;
  const max = -30;
  const percent = Math.max(
    0,
    Math.min(100, ((rssi - min) / (max - min)) * 100)
  );
  return Math.round(percent);
}

function getRsrpPercent(rsrp: number): number {
  if (!rsrp) return 0;
  // RSRP: -140dBm (0%) 到 -70dBm (100%)
  const min = -140;
  const max = -70;
  const percent = Math.max(
    0,
    Math.min(100, ((rsrp - min) / (max - min)) * 100)
  );
  return Math.round(percent);
}

function getRsrqPercent(rsrq: number): number {
  if (!rsrq) return 0;
  // RSRQ: -20dB (0%) 到 -3dB (100%)
  const min = -20;
  const max = -3;
  const percent = Math.max(
    0,
    Math.min(100, ((rsrq - min) / (max - min)) * 100)
  );
  return Math.round(percent);
}

function getSnrPercent(snr: number): number {
  if (!snr) return 0;
  // SNR: 0dB (0%) 到 30dB (100%)
  const min = 0;
  const max = 30;
  const percent = Math.max(0, Math.min(100, ((snr - min) / (max - min)) * 100));
  return Math.round(percent);
}

// 获取载波信息
function formatNrca(nrca: string, pre: string, type: number, index: number): string {
  if (!nrca) return '-';
  const carriers = nrca.split(';').filter(item => item.trim() !== '');
  const carrier = carriers[type];
  if (!carrier) return '-';
  const fields = carrier.split(',');
  // index 越界
  if (index < 0 || index >= fields.length) return '-';
  const value = fields[index];
  // 参数为空
  if (!value || value.trim() === '') return '-';
  const num = Number(value);
  if (Number.isNaN(num)) return '-';
  return pre + String(num);
}

// API 调用函数
async function callUbus<T>(request: any): Promise<T> {
  const response = await axios.post<UbusResponse<T>>('/api/ubus', request);
  if (response.data.code === 0) {
    return response.data.result as T;
  } else {
    throw new Error(response.data.msg || '接口返回错误');
  }
}

async function fetchAllData() {
  loading.value = true;
  error.value = null;

  try {
    // 并行请求所有数据
    const [
      netInfo,
      lan,
      wan,
      wan6,
      traffic,
      device,
      cpuTempData,
      simInfoData,
      simInfo2Data,
      wwanInfoData,
      lanUserData,
    ] = await Promise.all([
      callUbus<NetInfoResult>(netInfoRequest),
      callUbus<NetworkInterface>(lanRequest),
      callUbus<NetworkInterface>(wanRequest),
      callUbus<NetworkInterface>(wan6Request),
      callUbus<TrafficData>(trafficRequest),
      callUbus<DeviceInfo>(deviceInfoRequest),
      callUbus<CpuTemp>(cpuTempRequest),
      callUbus<SimInfo>(simInfoRequest),
      callUbus<SimInfo2>(simInfo2Request),
      callUbus<WwanInfo>(wwanRequest),
      callUbus<LanUserList>(lanUserListRequest),
    ]);

    data.value = netInfo;
    lanData.value = lan;
    wanData.value = wan;
    wan6Data.value = wan6;
    trafficData.value = traffic;
    deviceInfo.value = device;
    cpuTemp.value = cpuTempData;
    simInfo.value = simInfoData;
    simInfo2.value = simInfo2Data;
    wwanInfo.value = wwanInfoData;
    lanUserList.value = lanUserData;
  } catch (e: any) {
    error.value = e?.message || '请求失败';
    console.error('数据获取失败:', e);
  } finally {
    loading.value = false;
  }
}

function refresh() {
  fetchAllData().then((res) => {
    ElMessage.success('数据已刷新');
  });
}

function toggleAutoRefresh() {
  if (autoRefresh.value) {
    startAutoRefresh();
  } else {
    stopAutoRefresh();
  }
}

function updateRefreshInterval() {
  if (autoRefresh.value) {
    stopAutoRefresh();
    startAutoRefresh();
  }
}

function startAutoRefresh() {
  stopAutoRefresh();
  refreshTimer = window.setInterval(() => {
    fetchAllData();
  }, refreshInterval.value);
}

function stopAutoRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer);
    refreshTimer = null;
  }
}

// 一键ADB调试
function oneClickDebug() {
  callUbus<any>(openAdbRequest)
    .then(() => {
      ElMessage.success('已开启ADB调试模式');
    })
    .catch((err) => {
      ElMessage.error('请求失败：' + (err?.message || '未知错误'));
    });
}

function oneClickDebugClose() {
  callUbus<any>(closeAdbRequest)
    .then(() => {
      ElMessage.success('已关闭ADB调试模式');
    })
    .catch((err) => {
      ElMessage.error('请求失败：' + (err?.message || '未知错误'));
    });
}


// 短信转发
function smsForwardHandler() {
  ElNotification({
    title: '功能未实现',
    message: '短信转发功能尚未实现，敬请期待！',
    type: 'warning',
    duration: 5000,
  });
}

onMounted(() => {
  fetchAllData();
  if (autoRefresh.value) {
    startAutoRefresh();
  }
});

onUnmounted(() => {
  stopAutoRefresh();
});
</script>

<style scoped>
/* 基础样式 */
.page {
  color: white;
  min-height: 100vh;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 50%, #3b82f6 100%);
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 12px 16px;
  margin-bottom: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.title-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title {
  font-size: 28px;
  font-weight: 700;
  color: #ffffff;
  margin: 0;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.status-dot.online {
  background: #48bb78;
}

.status-dot.offline {
  background: #e53e3e;
}

.status-text {
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.8);
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* 控制区域 */
.controls {
  display: flex;
  align-items: center;
  gap: 16px;
}

.auto-refresh-controls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.9);
}

.checkbox-label input[type='checkbox'] {
  display: none;
}

.checkmark {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.6);
  border-radius: 4px;
  position: relative;
  transition: all 0.2s ease;
}

.checkbox-label input[type='checkbox']:checked + .checkmark {
  background: #4299e1;
  border-color: #4299e1;
}

.checkbox-label input[type='checkbox']:checked + .checkmark::after {
  content: '✓';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: white;
  font-size: 12px;
  font-weight: bold;
}

.auto-refresh-controls select {
  padding: 4px 6px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  font-size: 14px;
  color: rgba(255, 255, 255, 0.9);
  cursor: pointer;
  transition: border-color 0.2s ease;
}

.auto-refresh-controls select:focus {
  outline: none;
  border-color: #4299e1;
}

.auto-refresh-controls select:disabled {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.4);
  cursor: not-allowed;
}

/* 按钮样式 */
.btn {
  padding: 5px 10px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: #409eff5d;
  color: white;
  box-shadow: 0 4px 12px rgba(66, 153, 225, 0.1);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(66, 153, 225, 0.4);
}

.btn-danger {
  background: linear-gradient(135deg, #e53e3e, #c53030);
  color: white;
  box-shadow: 0 4px 12px rgba(229, 62, 62, 0.3);
}

.btn-danger:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(229, 62, 62, 0.4);
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top: 2px solid currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 状态页面 */
.loading,
.error,
.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 40px;
  text-align: center;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #e2e8f0;
  border-top: 4px solid #4299e1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 20px;
}

.loading p,
.error h3,
.error p,
.empty h3,
.empty p {
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
}

.error h3,
.empty h3 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 8px;
}

/* 内容区域 - 等宽网格 */
.content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(380px, 1fr));
  gap: 24px;
  align-items: stretch; /* 卡片等高 */
}

/* 卡片样式 */
.card {
  backdrop-filter: blur(20px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  overflow: hidden;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  display: flex;
  flex-direction: column;
}

.card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.hd {
  display: flex;
  gap: 10px;
}
.hd img {
  width: 24px;
}

.card-header {
  background: rgba(255, 255, 255, 0.1);
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #ffffff;
  margin: 0;
}

.card-content {
  padding: 16px;
  display: flex;
  flex: 1 1 auto;
  flex-direction: column;
}

/* 设备信息卡片 */
.device-info-card {
  grid-column: 1 / -1; /* 占据整行 */
}

/* 操作卡片 */
.device-actions-card {
  grid-column: 1 / -1; /* 占据整行 */
}

.device-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 20px;
}

.device-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.device-label {
  font-size: 12px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.7);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.device-values {
  display: flex;
  gap: 16px;
}

.load-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.load-label {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.6);
}

.load-value {
  font-size: 16px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
}

.memory-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.memory-bar {
  position: relative;
  height: 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.memory-fill {
  height: 100%;
  background: linear-gradient(
    90deg,
    #62718abb 0%,
    #68d391bb 70%,
    #63b3edbb 100%
  );
  border-radius: 10px;
  transition: width 0.3s ease;
}

.memory-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 11px;
  font-weight: 600;
  color: #ffffff;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
}

.memory-details {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 16px;
  color: rgba(255, 255, 255, 0.8);
}

.memory-used {
  font-weight: 600;
}

.memory-separator {
  color: rgba(255, 255, 255, 0.5);
}

.memory-total {
  color: rgba(255, 255, 255, 0.7);
}

.temp-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.temp-value {
  font-size: 18px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
}

.temp-bar {
  height: 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.temp-fill {
  height: 100%;
  background: linear-gradient(
    90deg,
    #62718abb 0%,
    #68d391bb 70%,
    #63b3edbb 100%
  );
  border-radius: 8px;
  transition: width 0.3s ease;
}

.uptime-value {
  font-size: 24px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

/* 信息网格 */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-item .label {
  font-size: 12px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.7);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-item .value {
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.9);
  word-break: break-all;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 标签样式 */
.tag {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.tag.success {
  background: rgba(54, 237, 69, 0.2);
  color: #75f655;
  border: 1px solid rgba(54, 237, 69, 0.3);
}

.tag.warning {
  background: rgba(237, 137, 54, 0.2);
  color: #f6ad55;
  border: 1px solid rgba(237, 137, 54, 0.3);
}

.tag.danger {
  background: rgba(229, 62, 62, 0.2);
  color: #fc8181;
  border: 1px solid rgba(229, 62, 62, 0.3);
}

/* 信号进度条 */
.signal-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 16px;
}

.signal-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.progress-bar {
  position: relative;
  height: 24px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.progress-fill {
  height: 100%;
  background: linear-gradient(
    90deg,
    #62718abb 0%,
    #68d391bb 50%,
    #63b3edbb 100%
  );
  border-radius: 12px;
  transition: width 0.3s ease;
  position: relative;
}

.progress-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 12px;
  font-weight: 600;
  color: #ffffff;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
  z-index: 1;
}

/* 顶部网络信息卡片中的信号条图标 */
.signal-bars {
  display: inline-flex;
  align-items: flex-end;
  gap: 3px;
}

.bar {
  width: 5px;
  height: 6px;
  background: rgba(255, 255, 255, 0.25);
  border-radius: 2px;
  transition: height 0.2s ease, background 0.2s ease;
}

.full-bar {
  height: 18px;
}

.bar:nth-child(1) {
  height: 6px;
}
.bar:nth-child(2) {
  height: 9px;
}
.bar:nth-child(3) {
  height: 12px;
}
.bar:nth-child(4) {
  height: 15px;
}
.bar:nth-child(5) {
  height: 18px;
}

.bar.active:nth-child(1) {
  background: #68d391;
}
.bar.active:nth-child(2) {
  background: #68d391;
}
.bar.active:nth-child(3) {
  background: #68d391;
}
.bar.active:nth-child(4) {
  background: #68d391;
}
.bar.active:nth-child(5) {
  background: #68d391;
}

.battery {
  position: relative;
  width: 30px;
  height: 14px;
  border: 2px solid #ffffffc4;
  border-radius: 3px;
  box-sizing: border-box;
  display: inline-block;
}

.battery-head {
  position: absolute;
  right: -4px;
  top: 1px;
  width: 3px;
  height: 8px;
  background: #ffffffc4;
  border-radius: 0 2px 2px 0;
}

.battery-level {
  height: 100%;
  background: #40d67a; /* 默认绿色 */
  transition: width 0.3s ease, background 0.3s ease;
  border-radius: 1px;
}

.battery-charging {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -55%) scale(1.8) rotate(20deg);
  font-size: 10px;
  color: #fff200;
  display: none; /* 默认不显示 */
}

/* 电量低时变红 */
.battery.low .battery-level {
  background: #f56565;
}

/* 充电状态显示闪电 */
.battery.charging .battery-charging {
  display: block;
}

/* 接口区域 */
.interface-section {
  margin-bottom: 24px;
}

.interface-section:last-child {
  margin-bottom: 0;
}

.interface-section h4 {
  font-size: 16px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  margin: 0 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 2px solid rgba(255, 255, 255, 0.2);
}

/* 网络接口状态网格布局 */
.interface-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
}

.info-grid-compact {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.info-grid-compact .info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-grid-compact .info-item .label {
  font-size: 11px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.6);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-grid-compact .info-item .value {
  font-size: 13px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.9);
  word-break: break-all;
}

/* 流量统计 */
.traffic-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 16px;
}

.traffic-item {
  text-align: center;
  padding: 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.traffic-label {
  font-size: 12px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.7);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}

.traffic-value {
  font-size: 18px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
}

.traffic-value.upload {
  color: #fc8181;
}

.traffic-value.download {
  color: #68d391;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page {
    padding: 12px;
  }

  .page-header {
    flex-direction: column;
    gap: 16px;
    padding: 16px;
  }

  .title {
    font-size: 24px;
  }

  .controls {
    flex-direction: column;
    gap: 12px;
    width: 100%;
  }

  .auto-refresh-controls {
    justify-content: center;
  }

  .content {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .device-stats {
    grid-template-columns: repeat(2, 1fr);
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .interface-grid {
    grid-template-columns: 1fr;
  }

  .info-grid-compact {
    grid-template-columns: 1fr;
  }

  .traffic-stats {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .page {
    padding: 8px;
  }

  .card-content {
    padding: 16px;
  }

  .device-stats {
    grid-template-columns: 1fr;
  }

  .traffic-item {
    padding: 12px;
  }

  .traffic-value {
    font-size: 16px;
  }
}

/* 暗色模式支持 */
@media (prefers-color-scheme: dark) {
  .page {
    color: white;
    background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
  }

  .page-header,
  .card,
  .loading,
  .error,
  .empty {
    background: rgba(15, 23, 42, 0.8);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .title,
  .card-header h3,
  .info-item .value,
  .traffic-value {
    color: #f1f5f9;
  }

  .status-text,
  .info-item .label,
  .traffic-label {
    color: #cbd5e1;
  }

  .card-header {
    background: rgba(30, 41, 59, 0.5);
  }

  .interface-section h4 {
    color: #f1f5f9;
    border-bottom-color: rgba(255, 255, 255, 0.1);
  }
}

.mytable{caption-side: bottom;}
.mytable{border-collapse: collapse;}
.mytable,tr,td{border: 1px solid rgb(59, 104, 141);font-size: 14px;text-align: center;}
.dbmstyle{color: rgb(104, 211, 145)}
</style>
