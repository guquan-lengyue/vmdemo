# VIRT-INSTALL(1) 虚拟化支持 VIRT-INSTALL(1)

## 名称（NAME）
virt-install - 创建新的虚拟机


## 语法（SYNOPSIS）
virt-install [选项]...


## 描述（DESCRIPTION）
virt-install 是一款命令行工具，可通过 libvirt 虚拟机监控程序（hypervisor）管理库创建新的 KVM（基于内核的虚拟机）、Xen 或 Linux 容器客户机（guest）。如需快速上手，请查看本文档末尾的“示例（EXAMPLES）”部分。

virt-install 工具支持通过 VNC 或 SPICE 等方式进行图形化安装，也支持通过串行控制台（serial console）进行文本模式安装。可配置客户机使用一个或多个虚拟磁盘、网络接口、音频设备、物理 USB 或 PCI 设备等。

安装介质可以是本地 ISO 文件或 CDROM 设备，也可以是通过 HTTP、FTP 远程托管的发行版安装树，或本地目录中的安装树。若使用安装树，virt-install 会仅获取启动安装过程所需的最小文件集，后续客户机可根据需要获取操作系统发行版的其余文件。同时支持 PXE 引导，以及导入现有磁盘镜像（从而跳过安装阶段）。

若提供合适的命令行参数，virt-install 可实现完全无人值守安装，客户机也能自行“ Kickstart（自动安装）”。这种方式便于实现客户机安装的自动化，可手动配置，或更简单地使用 `--unattended` 选项。

许多参数支持子选项，格式为 `opt1=foo,opt2=bar` 等。若需查看某个参数的完整子选项列表，可使用 `--option=?`，例如：`virt-install --disk=?`。

大多数选项为非必填项。若指定或检测到合适的 `--osinfo` 值，所有默认配置会自动填充并在终端输出中显示；否则，必填的最小选项集包括 `--memory`（内存）、客户机存储（`--disk` 或 `--filesystem`）以及安装方式选择。


## 连接到 libvirt（CONNECTING TO LIBVIRT）
### --connect
- **语法**：`--connect URI`
- **说明**：连接到非默认的虚拟机监控程序。若未指定，libvirt 会尝试选择最合适的默认监控程序。
- **有效选项示例**：
  - `qemu:///system`：用于创建由系统 libvirtd 实例运行的 KVM 和 QEMU 客户机，这是 virt-manager 使用的默认模式，也是大多数 KVM 用户所需的模式。
  - `qemu:///session`：用于创建由普通用户运行的 libvirtd 所管理的 KVM 和 QEMU 客户机。
  - `xen:///`：用于连接到 Xen 虚拟机监控程序。
  - `lxc:///`：用于创建 Linux 容器。


## 通用选项（GENERAL OPTIONS）
适用于所有类型客户机安装的通用配置参数。

### -n, --name
- **语法**：`-n, --name 名称`
- **说明**：新客户机虚拟机实例的名称。该名称在连接的虚拟机监控程序已知的所有客户机中必须唯一（包括当前未激活的客户机）。若需重新定义现有客户机，需先使用 virsh(1) 工具将其关闭（`virsh shutdown`）并删除（`virsh undefine`），再运行 virt-install。

### --memory
- **语法**：`--memory 选项`
- **说明**：为客户机分配的内存，单位为 MiB（兆字节）。此选项已弃用 `-r/--ram` 选项。
  - 支持子选项，如 `memory`、`currentMemory`、`maxMemory` 和 `maxMemory.slots`，这些子选项均映射到同名的 XML 配置值。
  - 为兼容旧版本，`memory` 映射到 `<currentMemory>` 元素，`maxmemory` 映射到 `<memory>` 元素。
  - 若需配置可热插拔的内存模块，请查看 `--memdev` 的说明。
- **查看子选项**：使用 `--memory=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsMemoryAllocation

### --memorybacking
- **语法**：`--memorybacking 选项`
- **说明**：此选项用于控制虚拟内存页如何由主机内存页提供支持。
- **查看子选项**：使用 `--memorybacking=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsMemoryBacking

### --arch
- **语法**：`--arch 架构`
- **说明**：为客户机虚拟机指定非主机原生的 CPU 架构。若省略，客户机将使用主机的 CPU 架构。

### --machine
- **语法**：`--machine 机器类型`
- **说明**：指定要模拟的机器类型。对于 Xen 或 KVM，通常无需指定此选项；但对于更特殊的架构，此选项可用于选择机器类型。

### --metadata
- **语法**：`--metadata 选项=值[,选项=值,...]`
- **说明**：为客户机指定元数据值，可选选项包括 `name`（名称）、`uuid`（唯一标识符）、`title`（标题）和 `description`（描述）。此选项已弃用 `-u/--uuid` 和 `--description`。
- **查看子选项**：使用 `--metadata=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsMetadata

### --events
- **语法**：`--events 选项=值[,选项=值,...]`
- **说明**：为客户机指定事件配置值，可选选项包括 `on_poweroff`（关机时）、`on_reboot`（重启时）和 `on_crash`（崩溃时）。
- **查看子选项**：使用 `--events=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsEvents

### --resource
- **语法**：`--resource 选项=值[,选项=值,...]`
- **说明**：为客户机指定资源分区配置。
- **查看子选项**：使用 `--resource=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#resPartition

### --sysinfo
- **语法**：`--sysinfo 选项=值[,选项=值,...]`
- **说明**：配置向虚拟机操作系统暴露的 sysinfo/SMBIOS（系统管理 BIOS）信息。示例：
  - `--sysinfo host`：特殊类型，将主机的 SMBIOS 信息暴露给虚拟机。
  - `--sysinfo emulate`：特殊类型，由虚拟机监控程序为虚拟机生成 SMBIOS 信息。
  - `--sysinfo bios.vendor=custom` 或 `--sysinfo smbios,bios.vendor=custom`：默认类型为 `smbios`，允许用户手动指定 SMBIOS 信息。
- **查看子选项**：使用 `--sysinfo=?` 可查看所有可用子选项。
- **详细文档**：SMBIOS XML 元素相关文档请参考 https://libvirt.org/formatdomain.html#elementsSysinfo 和 https://libvirt.org/formatdomain.html#elementsOSBIOS。

### --xml
- **语法**：`--xml 参数`
- **说明**：使用 XPath 语法直接编辑生成的 XML 配置。示例：
  ```bash
  virt-install --xml ./@foo=bar --xml ./newelement/subelement=1
  ```
  上述命令会将生成的 XML 修改为：
  ```xml
  <domain foo='bar' ...>
    ...
    <newelement>
      <subelement>1</subelement>
    </newelement>
  </domain>
  ```
- **子选项**：
  1. `--xml xpath.set=XPATH[=值]`：若未指定显式子选项，此为默认行为。格式为 `XPATH=值`（除非与 `xpath.value` 配合使用）。值的解释方式见下文。
  2. `--xml xpath.value=值`：`xpath.set` 仅被解释为 XPath 字符串，`xpath.value` 用作要设置的值。若需设置的字符串包含 `=` 符号，此选项可避免问题。若值为空，则视为删除该节点。
  3. `--xml xpath.create=XPATH`：创建空元素节点，适用于 `<readonly/>` 等布尔型元素。
  4. `--xml xpath.delete=XPATH`：删除 XPath 指定的节点及其所有子节点。

### xpath 子参数（xpath subarguments）
与 `--xml` 选项类似，大多数顶级选项都有 `xpath.*` 子选项。例如：`--disk xpath1.set=./@foo=bar,xpath2.create=./newelement` 会生成如下 XML 修改：
```xml
<disk foo="bar">
  <newelements/>
</disk>
```
当 virt-install 暂不支持某些设备的 XML 选项时，此功能可用于为设备单独设置 XML 配置。

### --qemu-commandline
- **语法**：`--qemu-commandline 参数`
- **说明**：将选项直接传递给 QEMU 模拟器，仅适用于 libvirt QEMU 驱动。此选项可接收参数字符串，示例：
  ```bash
  --qemu-commandline="-display gtk,gl=on"
  ```
  环境变量通过 `env` 指定，示例：
  ```bash
  --qemu-commandline=env=DISPLAY=:0.1
  ```
- **详细文档**：libvirt 该特性的详细说明请参考 https://libvirt.org/drvqemu.html#qemucommand。

### --vcpus
- **语法**：`--vcpus 选项`
- **说明**：为客户机配置的虚拟 CPU（vCPU）数量。若指定 `maxvcpus`，客户机运行时可热插拔 vCPU 至最大数量 `maxvcpus`，但启动时使用的 vCPU 数量为 `VCPUS`。
  - 还可指定 CPU 拓扑结构，包括 `sockets`（插槽数）、`dies`（芯片数）、`cores`（核心数）和 `threads`（线程数）。若部分值省略，其余值会自动填充，优先级为 `cores` > `sockets` > `threads`。优先选择核心数是因为这符合现代物理 CPU 的特性，更贴合客户机操作系统的预期。
  - `cpuset` 指定客户机可使用的物理 CPU，格式为逗号分隔的数字列表（支持范围或排除指定 CPU）。示例：
    - `0,2,3,5`：使用 CPU 0、2、3、5。
    - `1-5,^3,8`：使用 CPU 1、2、4、5、8（排除 3）。
  - 若传递值 `auto`，virt-install 会尝试利用 NUMA（非均匀内存访问）数据自动确定最优的 CPU 绑定（若可用）。
- **查看子选项**：使用 `--vcpus=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsCPUAllocation

### --numatune
- **语法**：`--numatune 选项`
- **说明**：调整客户机进程的 NUMA 策略。示例调用：
  ```bash
  --numatune 1,2,3,4-7
  --numatune 1-3,5,memory.mode=preferred
  ```
  指定内存分配的 NUMA 节点，语法与 `--vcpus cpuset=` 相同。`mode`（模式）可选值为 `interleave`（交错）、`preferred`（优先）或 `strict`（严格，默认值）。各模式详情请参考 `man 8 numactl`。
- **查看子选项**：使用 `--numatune=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsNUMATuning

### --memtune
- **语法**：`--memtune 选项`
- **说明**：调整客户机进程的内存策略。示例调用：
  ```bash
  --memtune 1000
  --memtune hard_limit=100,soft_limit=60,swap_hard_limit=150,min_guarantee=80
  ```
- **查看子选项**：使用 `--memtune=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsMemoryTuning

### --blkiotune
- **语法**：`--blkiotune 选项`
- **说明**：调整客户机进程的块 I/O（blkio）策略。示例调用：
  ```bash
  --blkiotune 100
  --blkiotune weight=100,device.path=/dev/sdc,device.weight=200
  ```
- **查看子选项**：使用 `--blkiotune=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsBlockTuning

### --cpu
- **语法**：`--cpu 型号[,+特性][,-特性][,match=匹配模式][,vendor=厂商],...`
- **说明**：配置向客户机暴露的 CPU 型号和 CPU 特性。唯一必填值为 `型号`（需是 libvirt 已知的有效 CPU 型号）。
  - libvirt 的特性策略值包括 `force`（强制启用）、`require`（要求启用）、`optional`（可选）、`disable`（禁用）、`forbid`（禁止），或简写形式 `+特性`（等同于 `force=特性`）和 `-特性`（等同于 `disable=特性`）。
  - 若指定精确的 CPU 型号，virt-install 会自动复制主机上可用的 CPU 特性，以缓解近期 CPU 推测执行侧信道漏洞和微架构存储缓冲区数据安全漏洞；但这会对性能产生一定影响，且可能导致无法迁移到未安装安全补丁的主机。可通过 `secure` 参数控制此行为，可选值为 `on`（启用，默认）和 `off`（禁用）。强烈建议保持启用，并确保所有虚拟化主机已安装最新的微码、内核和虚拟化软件。
- **示例**：
  1. `--cpu core2duo,+x2apic,disable=vmx`：暴露 core2duo CPU 型号，强制启用 x2apic，不暴露 vmx 特性。
  2. `--cpu host`：向客户机暴露主机的 CPU 配置，使客户机可利用主机 CPU 的多项特性（提升性能），但可能导致无法迁移到 CPU 型号不同的主机。
  3. `--cpu numa.cell0.memory=1234,numa.cell0.cpus=0-3,numa.cell1.memory=5678,numa.cell1.cpus=4-7`：指定两个 NUMA 节点，生成的 XML 如下：
     ```xml
     <cpu>
       <numa>
         <cell cpus="0-3" memory="1234"/>
         <cell cpus="4-7" memory="5678"/>
       </numa>
     </cpu>
     ```
  4. `--cpu host-passthrough,cache.mode=passthrough`：透传主机 CPU 的缓存信息。
- **查看子选项**：使用 `--cpu=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsCPU

### --cputune
- **语法**：`--cputune 选项`
- **说明**：调整客户机的 CPU 参数，可配置客户机 vCPU 绑定到主机的哪些物理 CPU。示例调用：
  ```bash
  --cputune vcpupin0.vcpu=0,vcpupin0.cpuset=0-3,vcpupin1.vcpu=1,vcpupin1.cpuset=4-7
  ```
- **查看子选项**：使用 `--cputune=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsCPUTuning

### --security, --seclabel
- **语法**：`--security, --seclabel type=类型[,label=标签][,relabel=yes|no],...`
- **说明**：配置客户机的安全标签（seclabel）域设置。`type`（类型）可选 `static`（静态）或 `dynamic`（动态）。`static` 配置需指定安全标签 `label`；若指定 `label` 但未指定 `type`，则默认使用 `static` 配置。
- **查看子选项**：使用 `--security=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#seclabel

### --keywrap
- **语法**：`--keywrap 选项`
- **说明**：指定客户机的 `<keywrap>` XML 配置，用于 S390 加密密钥管理操作。
- **查看子选项**：使用 `--keywrap=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#keywrap

### --iothreads
- **语法**：`--iothreads 选项`
- **说明**：指定客户机的 `<iothreads>` 和/或 `<iothreadids>` XML 配置。例如，若需配置 `<iothreads>4</iothreads>`，可使用 `--iothreads 4`。
- **查看子选项**：使用 `--iothreads=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsIOThreadsAllocation

### --features
- **语法**：`--features 特性=on|off,...`
- **说明**：设置客户机 `<features>` XML 中的元素启用（on）或禁用（off），可选特性包括 `acpi`（高级配置与电源接口）、`apic`（高级可编程中断控制器）、`eoi`（中断结束信号）、`privnet`（私有网络）和 `hyperv` 相关特性等。示例：
  1. `--features apic.eoi=on`：启用 APIC PV EOI（中断结束信号优化）。
  2. `--features hyperv.vapic.state=on,hyperv.spinlocks.state=off`：启用 Hyper-V VAPIC（虚拟高级可编程中断控制器），禁用自旋锁。
  3. `--features kvm.hidden.state=on`：向客户机隐藏 KVM 虚拟机监控程序标识。
  4. `--features pvspinlock=on`：通知客户机主机支持半虚拟化自旋锁（例如通过暴露 pvticketlocks 机制）。
  5. `--features gic.version=2`：仅适用于 ARM 架构，可选值为 `host`（主机版本）或具体版本号。
  6. `--features smm.state=on`：启用虚拟机监控程序的系统管理模式（SMM），部分 UEFI 固件可能要求此特性（QEMU 仅支持 q35 机器类型的 SMM）。
- **查看子选项**：使用 `--features=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsFeatures

### --clock
- **语法**：`--clock offset=偏移量,定时器选项=值,...`
- **说明**：配置客户机的 `<clock>` XML 配置，支持的选项包括：
  1. `--clock offset=偏移量`：设置时钟偏移量，例如 `utc`（世界协调时间）或 `localtime`（本地时间）。
  2. `--clock 定时器_present=no`：禁用布尔型定时器，`定时器` 可取值 `hpet`（高精度事件定时器）、`kvmclock`（KVM 时钟）等。
  3. `--clock 定时器_tickpolicy=值`：设置定时器的 tick 策略值，`定时器` 可取值 `rtc`（实时时钟）、`pit`（可编程间隔定时器）等，`值` 可取值 `catchup`（追赶）、`delay`（延迟）等。所有值的详情请参考 libvirt 文档。
- **查看子选项**：使用 `--clock=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsTime

### --pm
- **语法**：`--pm 选项`
- **说明**：配置客户机的电源管理特性。示例：
  ```bash
  --pm suspend_to_memi.enabled=on,suspend_to_disk.enabled=off
  ```
- **查看子选项**：使用 `--pm=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsPowerManagement

### --launchSecurity
- **语法**：`--launchSecurity 类型[,选项]`
- **说明**：为客户机启用启动安全机制（如 AMD SEV）。示例调用：
  1. 基础调用（使用默认策略 0x03，不提供 dhCert，无法与 SEV 固件交换数据）：
     ```bash
     --launchSecurity sev
     ```
  2. 指定策略（0x01，禁用调试，允许客户机密钥共享）：
     ```bash
     --launchSecurity sev,policy=0x01
     ```
  3. 提供会话信息（指定从 SEV 固件获取的会话 blob 和用于与 SEV 固件建立安全通信通道的 dhCert）：
     ```bash
     --launchSecurity sev,session=BASE64编码的会话字符串,dhCert=BASE64编码的dhCert字符串
     ```
  SEV 对 virtio 设备的使用有额外限制，如需查看包含 `--launchSecurity` 的完整 virt-install 调用示例，请参考“示例（EXAMPLES）”部分。
- **查看子选项**：使用 `--launchSecurity=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#launchSecurity


## 安装选项（INSTALLATION OPTIONS）
### -c, --cdrom
- **语法**：`--cdrom 路径`
- **说明**：用于虚拟机安装介质的 ISO 文件或 CDROM 设备。安装完成后，虚拟 CDROM 设备仍会挂载到虚拟机，但 ISO 或主机路径介质会被弹出。

### -l, --location
- **语法**：`-l, --location 选项`
- **说明**：发行版安装树的源路径。virt-install 可识别特定的发行版安装树，并获取可引导的内核/初始内存盘（initrd）对以启动安装。
  - `--location` 支持 `--extra-args`（内核参数）和 `--initrd-inject`（注入文件到 initrd）等选项。若需对 CDROM 介质使用这些选项，也可将 ISO 路径传递给 `--location`（但并非所有 CDROM 介质都支持）。
  - **路径格式**：
    1. `https://主机/路径`：包含可安装发行版镜像的 HTTP 服务器路径。
    2. `ftp://主机/路径`：包含可安装发行版镜像的 FTP 服务器路径。
    3. `ISO`：直接从 ISO 文件路径提取文件。
    4. `目录`：本地目录路径（包含可安装发行版镜像）。注意：初始引导后，客户机无法再访问该目录，因此操作系统安装程序需通过其他方式获取剩余安装介质。
  - **发行版示例路径**：
    - Fedora/Red Hat 系：`https://download.fedoraproject.org/pub/fedora/linux/releases/29/Server/x86_64/os`
    - Debian：`https://debian.osuosl.org/debian/dists/stable/main/installer-amd64/`
    - Ubuntu：`https://us.archive.ubuntu.com/ubuntu/dists/wily/main/installer-amd64/`
    - Suse：`https://download.opensuse.org/pub/opensuse/distribution/leap/42.3/repo/oss/`
  - **子选项**：`--location` 还支持 `kernel`（内核路径）和 `initrd`（initrd 路径）子选项，这些路径相对于指定的 `--location` URL/ISO，可用于在安装树中选择特定的内核/initrd 文件（适用于 virt-install/libosinfo 无法识别内核位置的场景）。示例：若有一个 libosinfo 无法识别的 ISO 文件 `my-unknown.iso`，其内核路径为 `kernel/fookernel`、initrd 路径为 `kernel/fooinitrd`，可使用以下命令：
    ```bash
    --location my-unknown.iso,kernel=kernel/fookernel,initrd=kernel/fooinitrd
    ```

### --pxe
- **说明**：通过 PXE（预启动执行环境）进行安装，仅指示虚拟机首次启动时从网络引导。

### --import
- **说明**：跳过操作系统安装过程，基于现有磁盘镜像创建客户机。用于引导的设备是通过 `--disk` 或 `--filesystem` 指定的第一个设备。

### -x, --extra-args
- **语法**：`-x, --extra-args 内核参数`
- **说明**：当通过 `--location` 安装客户机时，传递给安装程序的额外内核命令行参数。常见用途是指定 Anaconda Kickstart 自动安装文件，示例：
  ```bash
  --extra-args "ks=https://我的服务器/my.ks"
  ```

### --initrd-inject
- **语法**：`--initrd-inject 路径`
- **说明**：将指定路径的文件注入到通过 `--location` 获取的 initrd 根目录。此选项可用于无需网络托管 Kickstart 文件的自动安装，示例：
  ```bash
  --initrd-inject=/路径/to/my.ks --extra-args "ks=file:/my.ks"
  ```

### --install
- **说明**：用于各类安装操作的统一入口，支持多个子参数（类似 `--disk` 等选项）。此选项仅用于虚拟机安装操作，本质是配置首次启动。
  - **最简示例**：安装 Fedora 29：
    ```bash
    --install fedora29
    ```
    virt-install 会从 libosinfo 获取 `--location` URL，并自动填充默认配置。
  - **子选项**：
    1. `os=`：上述示例中的操作系统安装选项，显式写法为 `--install os=fedora29`。若未指定其他子选项，`os=` 为默认选项。
    2. `kernel=, initrd=`：指定用于安装介质的内核和 initrd 对。这两个文件会被复制到临时位置后再引导虚拟机，因此可与 `--initrd-inject` 配合使用，且不会修改源介质。若连接远程主机，介质会自动上传。示例：
       - 本地文件路径：`--install kernel=/路径/to/内核,initrd=/路径/to/initrd`
       - 网络路径（内核/initrd 会先下载到本地，再作为本地文件路径传递给虚拟机）：`--install kernel=https://127.0.0.1/安装树/内核,initrd=https://127.0.0.1/安装树/initrd`
       注意：此子选项仅用于安装阶段引导，若需设置虚拟机永久使用的内核，请使用 `--boot` 选项。
    3. `kernel_args=, kernel_args_overwrite=yes|no`：指定安装阶段的内核参数（对应 libvirt `<cmdline>` XML），可与 `kernel/initrd` 选项或 `--location` 介质配合使用。默认情况下，`kernel_args` 与 `--extra-args` 类似，会**追加**到 virt-install 为大多数 `--location` 安装设置的默认参数后；若需覆盖默认参数，需额外指定 `kernel_args_overwrite=yes`。
    4. `bootdev=`：指定安装阶段的引导设备（`hd` 硬盘、`cdrom` 光盘、`floppy` 软盘、`network` 网络），对应 libvirt `<os><boot dev=X>` XML。若需通过光盘或网络安装，直接使用 `--cdrom` 或 `--pxe` 更简单且兼容性更好；但此选项可在需要时对安装过程进行精细控制。
    5. `no_install=yes|no`：告知 virt-install 无需执行实际安装，仅创建虚拟机。`--import` 是此选项的别名，指定 `--boot` 但不指定其他安装选项也等同于设置此选项。已弃用的 `--live` 选项等同于 `--cdrom $ISO --install no_install=yes`。

### --reinstall DOMAIN
- **说明**：重新安装现有虚拟机，`DOMAIN` 可为主机名、UUID 或 ID 编号。virt-install 会从 libvirt 获取客户机 XML 配置，应用指定的安装配置修改，引导虚拟机执行安装过程，之后恢复为接近初始状态的 XML 配置。
  - 仅处理与安装相关的选项，`--name`、`--disk` 等其他虚拟机配置选项会被完全忽略。
  - 若 `--reinstall` 与 `--cdrom` 配合使用，会优先使用虚拟机已挂载的 CDROM 设备（若存在），否则会添加一个永久 CDROM 设备。

### --unattended
- **语法**：`--unattended [选项]`
- **说明**：通过 libosinfo 的自动安装脚本支持实现无人值守安装。libosinfo 包含多种发行版的自动安装脚本数据库（如 Red Hat Kickstart、Debian 安装脚本、Windows 无人值守安装脚本等）。最简调用方式是与 `--install` 配合：
  ```bash
  --install fedora29 --unattended
  ```
  Windows 安装示例：
  ```bash
  --cdrom /路径/to/我的/windows.iso --unattended
  ```
- **子选项**：
  1. `profile=`：选择 libosinfo 中的无人值守配置文件，大多数发行版提供 `desktop`（桌面版）和 `jeos`（精简版）配置文件。若未指定，virt-install 默认使用 `desktop`。
  2. `admin-password-file=`：用于设置虚拟机操作系统管理员/root 密码的文件，格式为 `admin-password-file=/路径/to/密码文件` 或 `admin-password-file=/dev/fd/n`（`n` 为密码文件的文件描述符）。注意：仅读取文件的第一行内容（包含空格，不包含换行符）。
  3. `user-login=`：虚拟机中的用户名，若未指定，默认使用当前主机用户名。注意：若以 `root` 用户运行 virt-install，必须指定此选项。
  4. `user-password-file=`：用于设置虚拟机用户密码的文件，格式与 `admin-password-file=` 相同。用户名由 `user-login` 指定或默认使用当前主机用户名。注意：仅读取文件的第一行内容（包含空格，不包含换行符）。
  5. `product-key=`：设置 Windows 产品密钥。

### --cloud-init
- **说明**：向虚拟机传递 cloud-init 元数据，会生成 cloud-init NoCloud ISO 文件并作为 CDROM 设备挂载到虚拟机（仅在首次启动时挂载）。此选项特别适用于发行版云镜像（默认锁定登录账户），可通过它初始化登录账户（如设置 root 密码）。
  - **最简调用**：仅使用 `--cloud-init`（无任何子选项），等同于 `--cloud-init root-password-generate=on,disable=on`，子选项说明见下文。
- **查看子选项**：使用 `--cloud-init=?` 可查看所有可用子选项。
- **子选项**：
  1. `root-password-generate=on`：为虚拟机生成新的 root 密码。启用此选项后，virt-install 会在控制台打印生成的密码，并暂停 10 秒供用户查看和复制。
  2. `disable=on`：在虚拟机后续启动时禁用 cloud-init。若不设置此选项，cloud-init 可能在每次启动时重置认证信息。
  3. `root-password-file=`：用于设置虚拟机 root 密码的文件，格式为 `root-password-file=/路径/to/密码文件` 或 `root-password-file=/dev/fd/n`（`n` 为密码文件的文件描述符）。注意：仅读取文件的第一行内容（包含空格，不包含换行符）。
  4. `meta-data=`：指定要直接添加到 ISO 的 cloud-init 元数据文件，此选项会忽略 `--cloud-init` 命令行中的其他元数据配置选项。
  5. `user-data=`：指定要直接添加到 ISO 的 cloud-init 用户数据文件，此选项会忽略 `--cloud-init` 命令行中的其他用户数据配置选项。
  6. `root-ssh-key=`：指定要注入到客户机的公钥，用于通过 SSH 登录 root 账户。示例：`root-ssh-key=/home/用户/.ssh/id_rsa.pub`。
  7. `clouduser-ssh-key`：指定要注入到客户机的公钥，用于通过 SSH 登录默认 cloud-init 用户账户（不同发行版云镜像的默认用户名不同，常见用户名参考 https://docs.openstack.org/image-guide/obtain-images.html）。
  8. `network-config=`：指定要直接添加到 ISO 的 cloud-init 网络配置文件。

### --boot
- **语法**：`--boot 启动选项`
- **说明**：可选，指定安装后虚拟机的启动配置。此选项可设置引导设备顺序、永久通过内核/initrd 引导（含可选内核参数），以及启用 BIOS 引导菜单（需 libvirt 0.8.3 及以上版本）。
  - `--boot` 可与其他安装选项（如 `--location`、`--cdrom` 等）配合使用，也可单独指定。单独指定时，行为类似 `--import` 安装选项：无“安装阶段”，直接按指定配置创建并启动客户机。
- **示例**：
  1. `--boot cdrom,fd,hd,network`：设置引导设备优先级为“光盘 → 软盘 → 硬盘 → 网络 PXE 引导”。
  2. `--boot kernel=内核文件,initrd=initrd文件,kernel_args="console=/dev/ttyS0"`：使客户机永久通过本地内核/initrd 对引导，并使用指定的内核参数。
  3. `--boot kernel=内核文件,initrd=initrd文件,dtb=设备树文件`：使客户机永久通过本地内核/initrd 对引导，并使用外部设备树二进制文件（DTB，部分非 x86 架构如 ARM 或 PPC 可能需要）。
  4. `--boot loader=BIOS路径`：使用指定路径的 BIOS 作为虚拟机 BIOS。
  5. `--boot bootmenu.enable=on,bios.useserial=on`：启用 BIOS 引导菜单，并通过串行控制台输出 BIOS 文本信息。
  6. `--boot init=初始化程序路径`：容器客户机的初始化程序路径。若已指定 `--filesystem` 根目录，virt-install 默认使用 `/sbin/init`；否则默认使用 `/bin/sh`。
  7. `--boot uefi`：配置虚拟机通过 UEFI 引导。为使 virt-install 识别正确的 UEFI 参数，libvirt 需通过 domcapabilities XML 暴露已知的 UEFI 二进制文件，因此此选项通常仅在使用正确配置的发行版软件包时生效。
  8. `--boot loader=/.../OVMF_CODE.fd,loader.readonly=yes,loader.type=pflash,nvram.template=/.../OVMF_VARS.fd,loader_secure=no`：指定虚拟机使用自定义 OVMF 二进制文件作为引导固件（映射为虚拟闪存芯片），并要求 libvirt 从自定义 `/.../OVMF_VARS.fd` 变量存储模板初始化虚拟机专用的 UEFI 变量存储。这是推荐的 UEFI 配置方式，若 `--boot uefi` 无法识别 UEFI 二进制文件，可使用此方式。若 UEFI 固件支持安全启动，可通过 `loader_secure` 启用。
- **查看子选项**：使用 `--boot=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsOS

### --idmap
- **语法**：`--idmap 选项`
- **说明**：若客户机配置声明了 UID（用户 ID）或 GID（组 ID）映射，会启用“用户命名空间”以应用这些映射。合适的 UID/GID 映射是容器在无 sVirt 隔离时实现安全的前提。
  - `--idmap` 可用于为 LXC 容器启用用户命名空间，示例：
    ```bash
    --idmap uid.start=0,uid.target=1000,uid.count=10,gid.start=0,gid.target=1000,gid.count=10
    ```
- **查看子选项**：使用 `--idmap=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsOSContainer


## 客户机操作系统选项（GUEST OS OPTIONS）
### --os-variant, --osinfo
- **语法**：`--osinfo [操作系统名称|选项1=值1,...]`
- **说明**：为特定操作系统优化客户机配置。大多数情况下，需指定操作系统或从安装介质检测操作系统，以启用 virtio 等性能关键特性。
  - **最简用法**：`--os-variant 操作系统名称` 或 `--osinfo 操作系统名称`，示例：`--osinfo fedora32`。
  - **子选项**：
    1. `name=, short-id=`：libosinfo 中的操作系统名称/短 ID，示例：`fedora32`、`win10`。
    2. `id=`：libosinfo 中的完整 URL 格式 ID，示例：`name=win10` 等同于 `id=http://microsoft.com/win/10`。
    3. `detect=on|off`：是否尝试从指定安装介质检测操作系统。目前仅对 URL 和 CDROM 安装尝试检测，且检测并非 100% 可靠。
    4. `require=on|off`：若设置为 `on`，则当未设置或未检测到操作系统值时，virt-install 会报错。
  - **示例**：
    1. `--osinfo detect=on,require=on`：尝试从安装介质检测操作系统，若检测失败则显式报错（确保 virt-install 不会回退到性能较差的配置）。
    2. `--osinfo detect=on,name=操作系统名称`：尝试从安装介质检测操作系统，若检测失败则使用指定的 `操作系统名称` 作为备用。
  - **默认行为**：若手动指定任何 `--osinfo` 值，其他设置默认关闭或未设置。virt-install 会为合适的安装介质自动尝试 `--osinfo detect=on`；若未检测到操作系统，大多数常见场景下会报错（此致命错误自 2022 年起添加）。可通过上述备用示例或禁用 `require` 选项解决此问题；若需快速恢复旧版非致命行为，可设置环境变量 `VIRTINSTALL_OSINFO_DISABLE_REQUIRE=1`。
  - **查看支持的操作系统**：使用 `virt-install --osinfo list` 可查看所有支持的操作系统变体；使用 `osinfo-query os` 可查看更完整的输出。
  - **注意**：`--os-variant` 和 `--osinfo` 互为别名，`--osinfo` 是推荐的新版命名方式。


## 存储选项（STORAGE OPTIONS）
### --disk
- **语法**：`--disk 选项`
- **说明**：指定客户机使用的存储介质，支持多种选项。磁盘字符串的通用格式为 `--disk 选项1=值1,选项2=值2,...`。
  - **最简调用**：创建 10GB 新磁盘镜像及关联磁盘设备：
    ```bash
    --disk size=10
    ```
    virt-install 会生成路径名并存储在虚拟机监控程序的默认镜像位置。
  - **指定介质的方式**：
    1. `--disk /某个/存储/路径[,选项1=值1]...`
    2. 显式指定以下参数之一：
       - `path`：存储介质的路径（现有或新建）。现有介质可以是文件或块设备；若路径不存在，意味着尝试创建新存储（需指定 `size` 值）。即使是远程主机，virt-install 也会尝试通过 libvirt 存储 API 自动创建指定路径。若虚拟机监控程序支持，`path` 也可以是网络 URL（如 `https://example.com/某个磁盘.img`），此时监控程序会直接访问存储，无需本地下载。
       - `pool`：用于创建新存储的现有 libvirt 存储池名称（需指定 `size` 值）。
       - `vol`：要使用的现有 libvirt 存储卷，格式为 `存储池名称/卷名称`。
  - **存储创建相关选项**：
    1. `size`：创建新存储时的大小（单位：GiB）。
    2. `sparse`：是否跳过新存储的完全分配，可选值 `yes`（是）或 `no`（否）。默认值为 `yes`（不完全分配），除非底层存储类型不支持。注意：完全分配客户机虚拟磁盘（`sparse=no`）的初始耗时通常可通过客户机内更快的安装速度抵消，因此建议使用此选项以确保稳定的高性能，并避免主机文件系统满导致客户机 I/O 错误。
    3. `format`：磁盘镜像格式。对于文件卷，可选格式包括 `raw`、`qcow2`、`vmdk` 等（所有可能值请参考 https://libvirt.org/storage.html 中的格式类型），此选项通常也映射到 `driver_type` 值。创建文件镜像时若未指定，默认格式为 `qcow2`；使用现有镜像时，此选项会覆盖 libvirt 的格式自动检测结果。
    4. `backing_store`：用作新创建镜像的后端存储的磁盘路径。
    5. `backing_format`：后端存储（backing_store）的磁盘镜像格式。
  - **设备配置子选项示例**：
    1. `device`：磁盘设备类型，可选值 `cdrom`（光盘）、`disk`（磁盘）、`lun`（逻辑单元号）或 `floppy`（软盘），默认值为 `disk`。
    2. `boot.order`：多磁盘客户机安装后需通过此参数确保正确引导，取值为 1、2、3...，值越小优先级越高（此选项也适用于其他可引导设备类型）。
    3. `target.bus` 或 `bus`：磁盘总线类型，可选值 `ide`、`sata`、`scsi`、`usb`、`virtio` 或 `xen`，默认值由虚拟机监控程序决定（并非所有监控程序都支持所有总线类型）。
    4. `readonly`：设置磁盘为只读，取值 `on`（是）或 `off`（否）。
    5. `shareable`：设置磁盘为可共享，取值 `on`（是）或 `off`（否）。
    6. `cache`：缓存模式，主机页缓存提供缓存内存。可选值：
       - `none`：无缓存。
       - `writethrough`：写穿透（提供读缓存）。
       - `directsync`：直接同步（绕过主机页缓存）。
       - `unsafe`：不安全（可能缓存所有内容并忽略客户机的刷新请求）。
       - `writeback`：写回（提供读写缓存）。
    7. `driver.discard`：是否将丢弃（又称“trim”或“unmap”）请求传递给文件系统，可选值 `unmap`（允许传递）或 `ignore`（忽略）（仅 QEMU 和 KVM 支持，需 libvirt 1.0.6 及以上版本）。
    8. `driver.name`：虚拟机监控程序访问指定存储时使用的驱动名称，通常无需用户设置。
    9. `driver.type`：虚拟机监控程序访问指定存储时使用的驱动格式/类型，通常无需用户设置。
    10. `driver.io`：磁盘 I/O 后端，可选值 `threads`（线程）、`native`（原生）或 `io_uring`。
    11. `driver.error_policy`：客户机遇到写错误时的处理策略，可选值 `stop`（停止）、`ignore`（忽略）或 `enospace`（空间不足）。
    12. `serial`：模拟磁盘设备的序列号，Linux 客户机中用于设置 `/dev/disk/by-id` 符号链接，示例序列号：`WD-WMAP9A966149`。
    13. `source.startupPolicy`：定义源文件不可访问时对磁盘的处理策略。
    14. `snapshot`：定义磁盘快照期间的默认行为。
  - **说明**：此选项已弃用 `-f/--file`、`-s/--file-size`、`--nonsparse` 和 `--nodisks`。
- **查看子选项**：使用 `--disk=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsDisks

### --filesystem
- **语法**：`--filesystem 选项`
- **说明**：指定要导出到客户机的主机目录。最简调用：
  ```bash
  --filesystem /主机上的源路径,/客户机中的目标路径
  ```
  此调用适用于新版 QEMU、Linux 客户机或 LXC 容器。对于 QEMU，目标路径仅为 sysfs 中的挂载提示，不会自动挂载。
- **子选项示例**：
  1. `type`：源目录类型，可选值 `mount`（挂载，默认）或 `template`（OpenVZ 模板）。
  2. `accessmode` 或 `mode`：客户机操作系统访问源目录的模式（仅 QEMU 和 `type=mount` 时使用），可选值 `mapped`（映射，默认）、`passthrough`（透传）或 `squash`（压缩）。详情请参考 libvirt 域 XML 文档。
  3. `source`：主机上要共享的目录。
  4. `target`：客户机中的挂载位置。
- **查看子选项**：使用 `--filesystem=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsFilesystems


## 网络选项（NETWORKING OPTIONS）
### -w, --network
- **语法**：`-w, --network 选项`
- **说明**：将客户机连接到主机网络，以下为指定网络类型的示例：
  1. `bridge=网桥名称`：连接到主机上名为“网桥名称”的网桥设备。若主机使用静态网络配置且客户机需与局域网进行完整的双向通信（入站/出站），或需对客户机进行实时迁移，建议使用此选项。
  2. `network=名称`：连接到主机上名为“名称”的虚拟网络。虚拟网络可通过 virsh 命令行工具列出、创建和删除；默认 libvirt 安装中通常存在一个名为 `default` 的虚拟网络。若主机使用动态网络（如 NetworkManager）或无线网络，建议使用虚拟网络，客户机将通过当前活跃的网络连接进行 NAT 转换以访问局域网。
  3. `type=direct,source=接口名[,source.mode=模式]`：通过 macvtap 直接连接到主机接口“接口名”。
  4. `user`：通过 SLIRP 连接到局域网，仅适用于以非特权用户身份运行 QEMU 客户机，提供有限的 NAT 功能。
  5. `none`：告知 virt-install 不添加任何默认网络接口。
  - **默认行为**：若省略 `--network`，会为客户机创建一个网卡。若主机存在连接物理接口的网桥设备，则使用该网桥；否则使用名为 `default` 的虚拟网络。此选项可多次指定以配置多个网卡。
- **子选项示例**：
  1. `model.type` 或 `model`：客户机识别到的网络设备型号，取值为虚拟机监控程序支持的任何网卡型号，例如 `e1000`、`rtl8139`、`virtio` 等。
  2. `mac.address` 或 `mac`：客户机的固定 MAC 地址。若省略或指定值 `RANDOM`，会自动生成合适的地址。注意：Xen 虚拟机要求 MAC 地址前 3 对字节为 `00:16:3e`；QEMU 或 KVM 虚拟机要求为 `52:54:00`。
  3. `filterref.filter`：控制 libvirt 中的防火墙和网络过滤，取值为 virsh `nwfilter` 子命令定义的任何网络过滤器。可用过滤器可通过 `virsh nwfilter-list` 查看，例如 `clean-traffic`（清洁流量）、`no-mac-spoofing`（禁止 MAC 欺骗）等。
  4. `virtualport.* 选项`：配置设备的虚拟端口配置文件，适用于 802.Qbg、802.Qbh、midonet 和 openvswitch 配置。
  - **说明**：此选项已弃用 `-m/--mac`、`-b/--bridge` 和 `--nonetworks`。
- **查看子选项**：使用 `--network=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsNICS


## 图形选项（GRAPHICS OPTIONS）
- **默认行为**：若未指定任何图形选项，virt-install 会在 `DISPLAY` 环境变量已设置时尝试选择合适的图形配置，否则使用 `--graphics none`。

### --graphics
- **语法**：`--graphics 类型,选项1=参数1,选项2=参数2,...`
- **说明**：指定图形显示配置，不涉及虚拟硬件配置，仅定义如何访问客户机的图形显示。通常无需用户指定此选项，virt-install 会自动选择实用的默认配置并启动相应的连接工具。
  - **通用格式**：`--graphics 类型,选项1=参数1,选项2=参数2,...`，示例：
    ```bash
    --graphics vnc,password=foobar
    ```
  - **支持的类型（TYPE）**：
    1. `vnc`：在客户机中设置虚拟控制台，并在主机中以 VNC 服务器形式暴露。若未指定 `port` 参数，VNC 服务器会在 5900 及以上端口中选择第一个可用端口。可通过 virsh 的 `vncdisplay` 命令获取分配的实际 VNC 显示（或使用 virt-viewer(1) 自动处理）。
    2. `spice`：通过 Spice 协议暴露客户机控制台，支持音频、USB 设备流传输及图形性能优化。使用 `spice` 类型等效于自动添加以下参数：`--video qxl --channel spicevmc`。
    3. `none`：不为客户机分配图形控制台。此时客户机需在第一个串行端口配置文本控制台（可通过 `--extra-args` 选项实现），并通过 `virsh console 名称` 命令连接到串行设备。
  - **支持的子选项**：
    1. `port`：为客户机控制台请求永久的静态端口号，适用于 `vnc` 和 `spice` 类型。
    2. `tlsPort`：指定 Spice TLS 端口，适用于 `spice` 类型。
    3. `websocket`：为客户机控制台请求 VNC WebSocket 端口，适用于 `vnc` 和 `spice` 类型；若指定 `-1`，则自动分配端口。
    4. `listen`：VNC/Spice 连接的监听地址，默认通常为 `127.0.0.1`（仅本地访问），部分虚拟机监控程序支持全局修改（例如 QEMU 驱动的默认值可在 `/etc/libvirt/qemu.conf` 中修改）。使用 `0.0.0.0` 允许其他机器访问；使用 `none` 表示显示服务器不监听任何端口（仅通过 libvirt Unix 套接字本地访问，如 `virt-viewer --attach`）；使用 `socket` 表示虚拟机在主机文件系统上监听 libvirt 生成的 Unix 套接字路径。
    5. `password`：请求控制台密码（连接时需输入）。注意：此信息可能会存入 virt-install 日志文件，请勿使用重要密码，适用于 `vnc` 和 `spice` 类型。
    6. `gl.enable`：是否使用 OpenGL 加速渲染，取值 `yes`（是）或 `no`（否），适用于 `spice` 类型。
    7. `gl.rendernode`：启用 `gl` 时使用的 DRM 渲染节点路径。
  - **说明**：此选项已弃用 `--vnc`、`--vncport`、`--vnclisten`、`-k/--keymap`、`--sdl` 和 `--nographics`。
- **查看子选项**：使用 `--graphics=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsGraphics

### --autoconsole
- **语法**：`--autoconsole 选项`
- **说明**：配置 virt-install 为虚拟机启动的交互式控制台。此选项非必填，默认行为会根据虚拟机配置自适应，但可通过此选项覆盖默认选择。
  - `--autoconsole graphical`：使用图形化工具 virt-viewer(1) 作为交互式控制台。
  - `--autoconsole text`：使用文本模式的 `virsh console` 作为交互式控制台。
  - `--autoconsole none`：等同于 `--noautoconsole`。
- **--noautoconsole**：不自动尝试连接到客户机控制台。
  - **注意**：指定此选项后，virt-install 会快速退出。若命令请求多阶段安装（如 `--cdrom` 或 `--location`），则安装阶段完成后，无论虚拟机是否请求重启，都会被关闭。若需虚拟机重启，需保持 virt-install 运行，可结合 `--wait` 选项（即使指定了 `--noautoconsole`）。


## 虚拟化选项（VIRTUALIZATION OPTIONS）
用于覆盖默认虚拟化类型选择的选项。

### -v, --hvm
- **说明**：请求使用全虚拟化（HVM），若主机同时支持半虚拟化和全虚拟化，则启用此模式。若连接到无硬件虚拟化支持的 Xen 虚拟机监控程序，此参数可能不可用。连接到基于 QEMU 的虚拟机监控程序时，此参数为默认隐含项。

### -p, --paravirt
- **说明**：将客户机配置为半虚拟化（paravirtualized）客户机。若主机同时支持半虚拟化和全虚拟化，且未指定 `--hvm` 或此参数，则默认使用半虚拟化。

### --container
- **说明**：将客户机配置为容器类型客户机。仅当虚拟机监控程序同时支持其他客户机类型时，才需指定此选项（例如 LXC 和 OpenVZ 默认为容器模式，此选项仅为兼容性保留）。

### --virt-type
- **说明**：指定安装客户机的虚拟机监控程序类型，可选值如 `kvm`、`qemu` 或 `xen`。可用选项可通过 `virsh capabilities` 命令在 `<domain>` 标签中查看。
  - **说明**：此选项已弃用 `--accelerate`（现已为默认行为）。若需安装纯 QEMU 客户机，使用 `--virt-type qemu`。


## 设备选项（DEVICE OPTIONS）
所有设备均支持一组 `address.*` 选项，用于配置设备在其父控制器或总线上的地址详情。参考文档：https://libvirt.org/formatdomain.html#elementsAddress。

### --controller
- **语法**：`--controller 选项`
- **说明**：为客户机附加控制器设备。
- **示例调用**：
  1. `--controller usb2`：添加完整的 USB2 控制器配置。
  2. `--controller usb3`：添加 USB3 控制器。
  3. `--controller type=usb,model=none`：完全禁用 USB。
  4. `--controller type=scsi,model=virtio-scsi`：添加 VirtIO SCSI 控制器。
  5. `--controller num_pcie_root_ports=数量`：若虚拟机默认使用 PCIe（PCI Express），控制默认添加的 PCIe 根端口控制器设备数量。
- **查看子选项**：使用 `--controller=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsControllers

### --input
- **语法**：`--input 选项`
- **说明**：为客户机附加输入设备，可选设备类型如鼠标（mouse）、平板（tablet）或键盘（keyboard）。
- **查看子选项**：使用 `--input=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsInput

### --hostdev, --host-device
- **语法**：`--hostdev, --host-device 选项`
- **说明**：将主机物理设备附加到客户机，以下为 `HOSTDEV` 的示例值：
  1. `--hostdev pci_0000_00_1b_0`：通过 libvirt 节点设备名称（可通过 `virsh nodedev-list` 查看）指定。
  2. `--hostdev 001.003`：通过总线和设备号（可通过 `lsusb` 查看）指定 USB 设备。
  3. `--hostdev 0x1234:0x5678`：通过厂商 ID 和产品 ID（可通过 `lsusb` 查看）指定 USB 设备。
  4. `--hostdev 1f.01.02`：通过 PCI 地址（可通过 `lspci` 查看）指定 PCI 设备。
  5. `--hostdev wlan0,type=net`：指定网络设备（适用于 LXC 容器）。
  6. `--hostdev /dev/net/tun,type=misc`：指定字符设备（适用于 LXC 容器）。
  7. `--hostdev /dev/sdf,type=storage`：指定块设备（适用于 LXC 容器）。
- **查看子选项**：使用 `--hostdev=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsHostDev

### --sound
- **语法**：`--sound 型号`
- **说明**：为客户机附加虚拟音频设备，`型号` 指定模拟声卡型号，可选值 `ich6`、`ich9`、`ac97`、`es1370`、`sb16`、`pcspk` 或 `default`。`default` 会尝试选择指定操作系统支持的最优型号。
  - **说明**：此选项已弃用旧版 `--soundhw` 选项。
- **查看子选项**：使用 `--sound=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsSound

### --audio
- **说明**：为客户机的 `--sound` 硬件配置主机音频输出。
- **查看子选项**：使用 `--audio=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#audio-backends

### --watchdog
- **语法**：`--watchdog 型号[,action=动作]`
- **说明**：为客户机附加虚拟硬件看门狗设备，需客户机内存在守护进程和设备驱动。当虚拟机看似挂起时，看门狗会触发信号；`action` 指定 libvirt 触发信号后的操作，可选值：
  - `reset`：强制重启客户机（默认）。
  - `poweroff`：强制关闭客户机电源。
  - `pause`：暂停客户机。
  - `none`：不执行任何操作。
  - `shutdown`：优雅关闭客户机（不推荐，因挂起的客户机可能无法响应优雅关闭）。
  - `型号`：模拟设备型号，可选 `i6300esb`（默认）或 `ib700`。
- **示例**：
  1. `--watchdog default`：使用推荐配置。
  2. `--watchdog i6300esb,action=poweroff`：使用 i6300esb 型号，触发时执行“关闭电源”操作。
- **查看子选项**：使用 `--watchdog=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsWatchdog

### --serial
- **语法**：`--serial 选项`
- **说明**：为客户机附加串行设备，支持多种选项。串行字符串的通用格式为 `--serial 类型,选项1=值1,选项2=值2,...`。`--serial` 和 `--parallel` 设备除特殊说明外，支持相同的选项。以下为字符设备重定向的类型示例：
  1. `--serial pty`：伪终端（PTY），分配的 PTY 会在运行中客户机的 XML 描述中列出。
  2. `--serial dev,path=主机路径`：主机设备，串行设备示例路径 `/dev/ttyS0`，并行设备示例路径 `/dev/parport0`。
  3. `--serial file,path=文件名`：将输出写入指定文件。
  4. `--serial tcp,host=主机:端口,source.mode=模式,protocol.type=协议`：TCP 网络控制台：
     - `模式`：`bind`（在 `主机:端口` 上等待连接）或 `connect`（将输出发送到 `主机:端口`），默认 `bind`。
     - `主机`：默认 `127.0.0.1`，`端口` 为必填项。
     - `协议`：`raw`（原始）或 `telnet`（默认 `raw`）；若为 `telnet`，端口会作为 telnet 服务器或客户端工作。
     - **示例**：
       - 在所有地址的 4567 端口等待连接：`--serial tcp,host=0.0.0.0:4567`
       - 连接到本地主机的 1234 端口：`--serial tcp,host=:1234,source.mode=connect`
       - 在本地主机的 2222 端口等待 telnet 连接（用户可通过 `telnet localhost 2222` 交互式连接到控制台）：`--serial tcp,host=:2222,source.mode=bind,source.protocol=telnet`
  5. `--serial udp,host=目标主机:端口,bind_host=绑定主机:绑定端口`：UDP 网络控制台：
     - `目标主机:端口`：输出发送的目标地址（默认 `目标主机` 为 `127.0.0.1`，`端口` 为必填项）。
     - `绑定主机:绑定端口`：可选的本地绑定地址（仅当指定 `绑定端口` 时，`绑定主机` 默认 `127.0.0.1`）。
     - **示例**：
       - 将输出发送到默认 syslog 端口（可能需编辑 `/etc/rsyslog.conf`）：`--serial udp,host=:514`
       - 将输出发送到远程主机 192.168.10.20 的 4444 端口（远程主机可通过 `nc -u -l 4444` 读取输出）：`--serial udp,host=192.168.10.20:4444`
  6. `--serial unix,path=Unix套接字路径,mode=模式`：Unix 套接字（参考 `man 7 unix`），`模式` 的行为和默认值与 `--serial tcp,mode=模式` 类似。
- **查看子选项**：使用 `--serial=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsCharSerial

### --parallel
- **语法**：`--parallel 选项`
- **说明**：指定并行设备，格式和选项与 `--serial` 基本相同。
- **查看子选项**：使用 `--parallel=?` 可查看所有可用子选项。
- **详细文档**：https://libvirt.org/formatdomain.html#elementsCharParallel

### --channel
- **语法**：`--channel 选项`
- **说明**：为客户机和主机之间添加通信通道设备，此选项使用与 `--serial` 和 `--parallel` 相同的选项指定通道的主机/源端，额外的 `target` 选项用于指定客户机如何识别通道。以下为字符设备重定向的类型示例：
  1. `--channel 源,target.type=guestfwd,target.address=主机:端口`：使用 QEMU 用户模式网络栈的通信通道，客户机可通过指定的 `主机:端口` 连接到通道。
