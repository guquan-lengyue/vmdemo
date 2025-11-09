# 系统ubuntu24.04
# 必要依赖
sudo apt install -y \
    virt-manager \
    qemu-kvm \
    libvirt-daemon-system \
    libvirt-clients \
    bridge-utils \
    virtinst \
    ovmf \ 
    virt-viewer
  
> + `virt-manager`  虚拟机图形化管理工具（GUI）
> + `qemu-kvm` KVM虚拟化底层硬件仿真核心
> + `libvirt-daemon-system`  虚拟化管理守护进程（系统级服务）
> + `libvirt-clients` `libvirt` 的命令行客户端工具集
> +  `bridge-utils` Linux 网络桥接配置工具
> + `virtinst` 虚拟机快速安装工具集（命令行）; 简化命令行创建虚拟机的流程
> +  `ovmf` KVM 虚拟机的 UEFI 固件（替代传统 BIOS）
> + `virt-viewer` 轻量级虚拟机图形控制台工具（专注 “显示”，不负责管理）


/backend 是go后端
启动命令 
```bash
# backend目录下
sudo go run main.go
```
/frontend 是vue前端
启动命令 
```
# 在frontend目录下
yarn 
yarn dev
```
