package service

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// DynamicAttributeService 处理动态属性获取的业务逻辑
type DynamicAttributeService struct {
	cache     map[string]map[string]interface{} // 缓存结构：域名 -> 属性名 -> 属性值
	cacheLock sync.RWMutex
	cacheTTL  time.Duration // 缓存有效期
}

// NewDynamicAttributeService 创建一个新的DynamicAttributeService实例
func NewDynamicAttributeService() *DynamicAttributeService {
	return &DynamicAttributeService{
		cache:    make(map[string]map[string]interface{}),
		cacheTTL: 24 * time.Hour, // 默认缓存24小时
	}
}

// GetDynamicAttributes 获取域名的所有动态属性
func (s *DynamicAttributeService) GetDynamicAttributes(domain string) (map[string]interface{}, error) {
	// 检查缓存
	if attrs := s.getFromCache(domain); attrs != nil {
		return attrs, nil
	}

	// 创建结果map
	result := make(map[string]interface{})

	// 并发获取各种动态属性
	var wg sync.WaitGroup
	var mu sync.Mutex // 用于保护result map的并发写入
	errChan := make(chan error, 5)

	// 获取WHOIS信息
	wg.Add(1)
	go func() {
		defer wg.Done()
		whoisInfo, err := s.getWhoisInfo(domain)
		if err != nil {
			errChan <- fmt.Errorf("获取WHOIS信息失败: %w", err)
			return
		}
		mu.Lock()
		for k, v := range whoisInfo {
			result[k] = v
		}
		mu.Unlock()
	}()

	// 获取Alexa排名
	wg.Add(1)
	go func() {
		defer wg.Done()
		rank, err := s.getAlexaRank(domain)
		if err != nil {
			errChan <- fmt.Errorf("获取Alexa排名失败: %w", err)
			return
		}
		mu.Lock()
		result["alexa_rank"] = rank
		mu.Unlock()
	}()

	// 获取搜索量
	wg.Add(1)
	go func() {
		defer wg.Done()
		searchVolume, err := s.getSearchVolume(domain)
		if err != nil {
			errChan <- fmt.Errorf("获取搜索量失败: %w", err)
			return
		}
		mu.Lock()
		result["search_volume"] = searchVolume
		mu.Unlock()
	}()

	// 获取相关域名状态
	wg.Add(1)
	go func() {
		defer wg.Done()
		relatedDomains, err := s.getRelatedDomainsStatus(domain)
		if err != nil {
			errChan <- fmt.Errorf("获取相关域名状态失败: %w", err)
			return
		}
		mu.Lock()
		for k, v := range relatedDomains {
			result[k] = v
		}
		mu.Unlock()
	}()

	// 获取社交媒体和电商数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		socialData, err := s.getSocialAndEcommerceData(domain)
		if err != nil {
			errChan <- fmt.Errorf("获取社交媒体和电商数据失败: %w", err)
			return
		}
		mu.Lock()
		for k, v := range socialData {
			result[k] = v
		}
		mu.Unlock()
	}()

	// 等待所有goroutine完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误
	var errs []string
	for err := range errChan {
		errs = append(errs, err.Error())
	}

	// 如果有错误，但仍然获取了一些数据，我们继续处理
	// 只有在完全没有数据的情况下才返回错误
	if len(errs) > 0 && len(result) == 0 {
		return nil, fmt.Errorf("获取动态属性失败: %s", strings.Join(errs, "; "))
	}

	// 保存到缓存
	s.saveToCache(domain, result)

	return result, nil
}

// getFromCache 从缓存中获取域名属性
func (s *DynamicAttributeService) getFromCache(domain string) map[string]interface{} {
	s.cacheLock.RLock()
	defer s.cacheLock.RUnlock()

	if attrs, ok := s.cache[domain]; ok {
		return attrs
	}
	return nil
}

// saveToCache 保存域名属性到缓存
func (s *DynamicAttributeService) saveToCache(domain string, attrs map[string]interface{}) {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	s.cache[domain] = attrs

	// 设置过期时间（可以使用一个goroutine在后台清理过期缓存）
	go func() {
		time.Sleep(s.cacheTTL)
		s.cacheLock.Lock()
		delete(s.cache, domain)
		s.cacheLock.Unlock()
	}()
}

// getWhoisInfo 获取域名的WHOIS信息
func (s *DynamicAttributeService) getWhoisInfo(domain string) (map[string]interface{}, error) {
	// 在实际系统中，这里应该调用WHOIS API或解析WHOIS服务器响应
	// 这里使用模拟数据
	result := make(map[string]interface{})

	// 模拟数据：根据域名生成一些随机但看起来合理的注册信息
	domainHash := int(time.Now().Unix() % 10)
	if strings.Contains(domain, "a") {
		domainHash = (domainHash + 1) % 10
	}
	if strings.Contains(domain, "e") {
		domainHash = (domainHash + 2) % 10
	}

	// 注册日期：1-10年前
	registerYearsAgo := 1 + domainHash
	registerDate := time.Now().AddDate(-registerYearsAgo, 0, 0)
	result["register_date"] = registerDate.Format("2006-01-02")

	// 到期日期：1-3年后
	expireYearsLater := 1 + (domainHash % 3)
	expireDate := time.Now().AddDate(expireYearsLater, 0, 0)
	result["expire_date"] = expireDate.Format("2006-01-02")

	// 注册商
	registrars := []string{"GoDaddy", "Namecheap", "Alibaba Cloud", "Tencent Cloud", "NameSilo"}
	result["registrar"] = registrars[domainHash%len(registrars)]

	return result, nil
}

// getAlexaRank 获取域名的Alexa排名
func (s *DynamicAttributeService) getAlexaRank(domain string) (int, error) {
	// 在实际系统中，这里应该调用Alexa API
	// 这里使用模拟数据

	// 根据域名长度和字符生成一个看起来合理的排名
	// 短域名和含有常见词的域名排名更高
	rank := 1000000 // 默认排名

	// 域名越短，排名越高
	if len(domain) < 10 {
		rank = rank / (12 - len(domain))
	}

	// 含有常见词的域名排名更高
	commonWords := []string{"news", "shop", "blog", "tech", "game", "app", "web", "cloud"}
	for _, word := range commonWords {
		if strings.Contains(domain, word) {
			rank = rank / 2
			break
		}
	}

	// 添加一些随机性
	rank = rank + int(time.Now().Unix()%10000) - 5000
	if rank < 100 {
		rank = 100 + int(time.Now().Unix()%900)
	}

	return rank, nil
}

// getSearchVolume 获取域名相关关键词的搜索量
func (s *DynamicAttributeService) getSearchVolume(domain string) (int, error) {
	// 在实际系统中，这里应该调用搜索API，如Google Keyword Planner或百度指数
	// 这里使用模拟数据

	// 提取域名主体（不含TLD）
	parts := strings.Split(domain, ".")
	keyword := parts[0]

	// 根据关键词长度和字符生成一个看起来合理的搜索量
	volume := 1000 // 默认搜索量

	// 关键词越短，搜索量越高
	if len(keyword) < 6 {
		volume = volume * (7 - len(keyword))
	}

	// 含有常见词的关键词搜索量更高
	commonWords := []string{"news", "shop", "blog", "tech", "game", "app", "web", "cloud"}
	for _, word := range commonWords {
		if strings.Contains(keyword, word) {
			volume = volume * 3
			break
		}
	}

	// 添加一些随机性
	volume = volume + int(time.Now().Unix()%1000) - 500
	if volume < 100 {
		volume = 100 + int(time.Now().Unix()%900)
	}

	return volume, nil
}

// getRelatedDomainsStatus 获取相关域名的注册状态
func (s *DynamicAttributeService) getRelatedDomainsStatus(domain string) (map[string]interface{}, error) {
	// 在实际系统中，这里应该调用WHOIS API查询相关域名
	// 这里使用模拟数据
	result := make(map[string]interface{})

	// 提取域名主体（不含TLD）
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("无效的域名格式: %s", domain)
	}

	domainName := parts[0]
	mainTLD := parts[len(parts)-1]

	// 检查相关TLD
	relatedTLDs := []string{"com", "net", "org", "co", "cc", "io", "ai"}
	for _, tld := range relatedTLDs {
		if tld == mainTLD {
			continue // 跳过主域名的TLD
		}

		// 构建相关域名（用于日志或调试）
		_ = fmt.Sprintf("%s.%s", domainName, tld)

		// 模拟注册状态：根据域名和TLD生成一个看起来合理的状态
		// 越常见的TLD，被注册的可能性越高
		isRegistered := false
		if tld == "com" || tld == "net" || tld == "org" {
			isRegistered = (len(domainName) < 6) // 短域名在常见TLD下更可能被注册
		} else {
			isRegistered = (len(domainName) < 4) // 非常短的域名在其他TLD下更可能被注册
		}

		// 添加一些随机性
		if time.Now().Unix()%2 == 0 {
			isRegistered = !isRegistered
		}

		status := "未注册"
		if isRegistered {
			registerYearsAgo := 1 + int(time.Now().Unix()%5)
			registerDate := time.Now().AddDate(-registerYearsAgo, 0, 0)
			status = fmt.Sprintf("在 %s 注册", registerDate.Format("2006.01.02"))
		}

		result[fmt.Sprintf("related_domain_%s", tld)] = status
	}

	return result, nil
}

// getSocialAndEcommerceData 获取社交媒体和电商数据
func (s *DynamicAttributeService) getSocialAndEcommerceData(domain string) (map[string]interface{}, error) {
	// 在实际系统中，这里应该调用各种API获取社交媒体和电商数据
	// 这里使用模拟数据
	result := make(map[string]interface{})

	// 提取域名主体（不含TLD）
	parts := strings.Split(domain, ".")
	keyword := parts[0]

	// 贴吧数量
	tiebaPosts := 1000 + len(keyword)*500 + int(time.Now().Unix()%10000)
	result["tieba_posts"] = tiebaPosts

	// 百科系数
	baikeIndex := 1000 + len(keyword)*200 + int(time.Now().Unix()%5000)
	result["baike_index"] = baikeIndex

	// 词典记录
	hasDictRecord := len(keyword) < 6 || time.Now().Unix()%2 == 0
	result["dict_record"] = hasDictRecord

	// 360搜索指数
	search360Index := 500 + len(keyword)*100 + int(time.Now().Unix()%2000)
	result["search_360_index"] = search360Index

	// 传媒系数
	mediaIndex := 10000 + len(keyword)*1000 + int(time.Now().Unix()%100000)
	result["media_index"] = mediaIndex

	// 社交系数
	socialIndex := 5000 + len(keyword)*500 + int(time.Now().Unix()%50000)
	result["social_index"] = socialIndex

	// 淘宝商品数量
	taobaoProducts := 100 + len(keyword)*50 + int(time.Now().Unix()%1000)
	result["taobao_products"] = taobaoProducts

	return result, nil
}
