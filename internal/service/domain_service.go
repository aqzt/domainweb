package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"domainweb/internal/model"
	"domainweb/internal/repository"
)

// DomainService 处理域名估价的业务逻辑
type DomainService struct {
	repo               *repository.DomainRepository
	dynamicAttrService *DynamicAttributeService
}

// NewDomainService 创建一个新的DomainService实例
func NewDomainService(repo *repository.DomainRepository) *DomainService {
	return &DomainService{
		repo:               repo,
		dynamicAttrService: NewDynamicAttributeService(),
	}
}

// EstimateDomain 估算域名价值和品相等级
func (s *DomainService) EstimateDomain(domainName string) (*model.EstimationResult, error) {
	// 基础价格和等级
	const basePrice = 25.0
	const baseGrade = -0.5

	// 解析域名
	domain, err := s.parseDomain(domainName)
	if err != nil {
		return nil, fmt.Errorf("解析域名失败: %w", err)
	}

	// 获取基础属性
	baseAttributes, err := s.repo.GetAttributesByType("基础属性")
	if err != nil {
		return nil, fmt.Errorf("获取基础属性失败: %w", err)
	}

	// 获取其他属性
	otherAttributes, err := s.repo.GetAttributesByType("其他属性")
	if err != nil {
		return nil, fmt.Errorf("获取其他属性失败: %w", err)
	}

	// 计算基础属性的影响
	var baseAttrDetails []model.AttributeDetail
	totalPriceFactor := 1.0
	totalGradeFactor := 0.0

	// 处理TLD属性
	tldAttrs, err := s.repo.GetTLDAttributes()
	if err != nil {
		return nil, fmt.Errorf("获取TLD属性失败: %w", err)
	}

	if tldAttr, ok := tldAttrs[domain.TLD]; ok {
		totalPriceFactor *= tldAttr.PriceFactor
		totalGradeFactor += tldAttr.GradeFactor
		baseAttrDetails = append(baseAttrDetails, model.AttributeDetail{
			Name:        tldAttr.AttributeName,
			Value:       domain.TLD,
			Description: fmt.Sprintf("%s后缀", domain.TLD),
			PriceFactor: tldAttr.PriceFactor,
			GradeFactor: tldAttr.GradeFactor,
		})
	}

	// 处理长度属性
	for _, attr := range baseAttributes {
		if strings.Contains(attr.AttributeName, "位长度") &&
			fmt.Sprintf("%d", domain.Length) == attr.AttributeValue {
			totalPriceFactor *= attr.PriceFactor
			totalGradeFactor += attr.GradeFactor
			baseAttrDetails = append(baseAttrDetails, model.AttributeDetail{
				Name:        attr.AttributeName,
				Value:       fmt.Sprintf("%d", domain.Length),
				Description: fmt.Sprintf("%d位长度", domain.Length),
				PriceFactor: attr.PriceFactor,
				GradeFactor: attr.GradeFactor,
			})
			break
		}
	}

	// 处理结构属性
	for _, attr := range baseAttributes {
		if strings.Contains(attr.AttributeName, "结构") &&
			domain.Structure == attr.AttributeValue {
			totalPriceFactor *= attr.PriceFactor
			totalGradeFactor += attr.GradeFactor
			baseAttrDetails = append(baseAttrDetails, model.AttributeDetail{
				Name:        attr.AttributeName,
				Value:       domain.Structure,
				Description: fmt.Sprintf("%s结构", domain.Structure),
				PriceFactor: attr.PriceFactor,
				GradeFactor: attr.GradeFactor,
			})
			break
		}
	}

	// 获取动态属性
	var otherAttrDetails []model.AttributeDetail

	// 调用动态属性服务获取实时数据
	dynamicAttrs, err := s.dynamicAttrService.GetDynamicAttributes(domainName)
	if err != nil {
		// 如果获取动态属性失败，记录错误但继续处理
		fmt.Printf("获取动态属性失败: %v\n", err)
	}

	// 处理动态属性
	if dynamicAttrs != nil {
		// 处理Alexa排名
		if alexaRank, ok := dynamicAttrs["alexa_rank"].(int); ok {
			var alexaAttr model.DomainAttribute

			// 根据Alexa排名设置不同的影响因子
			if alexaRank < 10000 {
				alexaAttr = model.DomainAttribute{
					AttributeName:  "Alexa排名优秀",
					AttributeType:  "其他属性",
					PriceFactor:    2.5,
					GradeFactor:    0.8,
					AttributeValue: strconv.Itoa(alexaRank),
				}
			} else if alexaRank < 100000 {
				alexaAttr = model.DomainAttribute{
					AttributeName:  "Alexa排名良好",
					AttributeType:  "其他属性",
					PriceFactor:    1.8,
					GradeFactor:    0.5,
					AttributeValue: strconv.Itoa(alexaRank),
				}
			} else if alexaRank < 1000000 {
				alexaAttr = model.DomainAttribute{
					AttributeName:  "Alexa排名一般",
					AttributeType:  "其他属性",
					PriceFactor:    1.2,
					GradeFactor:    0.2,
					AttributeValue: strconv.Itoa(alexaRank),
				}
			} else {
				alexaAttr = model.DomainAttribute{
					AttributeName:  "Alexa排名",
					AttributeType:  "其他属性",
					PriceFactor:    1.0,
					GradeFactor:    0.0,
					AttributeValue: strconv.Itoa(alexaRank),
				}
			}

			totalPriceFactor *= alexaAttr.PriceFactor
			totalGradeFactor += alexaAttr.GradeFactor
			otherAttrDetails = append(otherAttrDetails, model.AttributeDetail{
				Name:        alexaAttr.AttributeName,
				Value:       alexaAttr.AttributeValue,
				Description: fmt.Sprintf("Alexa 排名 %s", alexaAttr.AttributeValue),
				PriceFactor: alexaAttr.PriceFactor,
				GradeFactor: alexaAttr.GradeFactor,
			})
		}

		// 处理搜索量
		if searchVolume, ok := dynamicAttrs["search_volume"].(int); ok {
			var searchAttr model.DomainAttribute

			// 根据搜索量设置不同的影响因子
			if searchVolume > 10000 {
				searchAttr = model.DomainAttribute{
					AttributeName:  "搜索量巨大",
					AttributeType:  "其他属性",
					PriceFactor:    3.0,
					GradeFactor:    0.9,
					AttributeValue: strconv.Itoa(searchVolume),
				}
			} else if searchVolume > 5000 {
				searchAttr = model.DomainAttribute{
					AttributeName:  "搜索量很高",
					AttributeType:  "其他属性",
					PriceFactor:    2.2,
					GradeFactor:    0.7,
					AttributeValue: strconv.Itoa(searchVolume),
				}
			} else if searchVolume > 1000 {
				searchAttr = model.DomainAttribute{
					AttributeName:  "搜索量较高",
					AttributeType:  "其他属性",
					PriceFactor:    1.8,
					GradeFactor:    0.6,
					AttributeValue: strconv.Itoa(searchVolume),
				}
			} else {
				searchAttr = model.DomainAttribute{
					AttributeName:  "搜索量",
					AttributeType:  "其他属性",
					PriceFactor:    1.0,
					GradeFactor:    0.0,
					AttributeValue: strconv.Itoa(searchVolume),
				}
			}

			totalPriceFactor *= searchAttr.PriceFactor
			totalGradeFactor += searchAttr.GradeFactor
			otherAttrDetails = append(otherAttrDetails, model.AttributeDetail{
				Name:        searchAttr.AttributeName,
				Value:       searchAttr.AttributeValue,
				Description: fmt.Sprintf("搜索量 %s", searchAttr.AttributeValue),
				PriceFactor: searchAttr.PriceFactor,
				GradeFactor: searchAttr.GradeFactor,
			})
		}

		// 处理相关域名状态
		for key, value := range dynamicAttrs {
			if strings.HasPrefix(key, "related_domain_") {
				tld := strings.TrimPrefix(key, "related_domain_")
				status, ok := value.(string)
				if !ok {
					continue
				}

				var relatedAttr model.DomainAttribute
				isRegistered := !strings.Contains(status, "未注册")

				if isRegistered {
					// 相关域名已注册，降低当前域名价值
					relatedAttr = model.DomainAttribute{
						AttributeName:  fmt.Sprintf("%s相关域名已注册", tld),
						AttributeType:  "其他属性",
						PriceFactor:    0.86,
						GradeFactor:    -0.1,
						AttributeValue: status,
					}
				} else {
					// 相关域名未注册，进一步降低当前域名价值
					relatedAttr = model.DomainAttribute{
						AttributeName:  fmt.Sprintf("%s相关域名未注册", tld),
						AttributeType:  "其他属性",
						PriceFactor:    0.65,
						GradeFactor:    -0.2,
						AttributeValue: "未注册",
					}
				}

				totalPriceFactor *= relatedAttr.PriceFactor
				totalGradeFactor += relatedAttr.GradeFactor
				otherAttrDetails = append(otherAttrDetails, model.AttributeDetail{
					Name:        relatedAttr.AttributeName,
					Value:       relatedAttr.AttributeValue,
					Description: fmt.Sprintf("%s.%s %s", strings.Split(domainName, ".")[0], tld, relatedAttr.AttributeValue),
					PriceFactor: relatedAttr.PriceFactor,
					GradeFactor: relatedAttr.GradeFactor,
				})
			}
		}

		// 处理社交媒体和电商数据
		if tiebaPosts, ok := dynamicAttrs["tieba_posts"].(int); ok && tiebaPosts > 0 {
			var tiebaAttr model.DomainAttribute

			if tiebaPosts > 10000 {
				tiebaAttr = model.DomainAttribute{
					AttributeName:  "贴吧数量巨大",
					AttributeType:  "其他属性",
					PriceFactor:    2.25,
					GradeFactor:    0.6,
					AttributeValue: strconv.Itoa(tiebaPosts),
				}
			} else {
				tiebaAttr = model.DomainAttribute{
					AttributeName:  "贴吧数量",
					AttributeType:  "其他属性",
					PriceFactor:    1.5,
					GradeFactor:    0.3,
					AttributeValue: strconv.Itoa(tiebaPosts),
				}
			}

			totalPriceFactor *= tiebaAttr.PriceFactor
			totalGradeFactor += tiebaAttr.GradeFactor
			otherAttrDetails = append(otherAttrDetails, model.AttributeDetail{
				Name:        tiebaAttr.AttributeName,
				Value:       tiebaAttr.AttributeValue,
				Description: fmt.Sprintf("贴吧数量 %s", tiebaAttr.AttributeValue),
				PriceFactor: tiebaAttr.PriceFactor,
				GradeFactor: tiebaAttr.GradeFactor,
			})
		}

		// 处理淘宝商品数量
		if taobaoProducts, ok := dynamicAttrs["taobao_products"].(int); ok && taobaoProducts > 0 {
			var taobaoAttr model.DomainAttribute

			if taobaoProducts > 1000 {
				taobaoAttr = model.DomainAttribute{
					AttributeName:  "淘宝商品数量多",
					AttributeType:  "其他属性",
					PriceFactor:    1.5,
					GradeFactor:    0.3,
					AttributeValue: strconv.Itoa(taobaoProducts),
				}
			} else {
				taobaoAttr = model.DomainAttribute{
					AttributeName:  "淘宝商品",
					AttributeType:  "其他属性",
					PriceFactor:    1.18,
					GradeFactor:    0.1,
					AttributeValue: strconv.Itoa(taobaoProducts),
				}
			}

			totalPriceFactor *= taobaoAttr.PriceFactor
			totalGradeFactor += taobaoAttr.GradeFactor
			otherAttrDetails = append(otherAttrDetails, model.AttributeDetail{
				Name:        taobaoAttr.AttributeName,
				Value:       taobaoAttr.AttributeValue,
				Description: fmt.Sprintf("淘宝商品 %s", taobaoAttr.AttributeValue),
				PriceFactor: taobaoAttr.PriceFactor,
				GradeFactor: taobaoAttr.GradeFactor,
			})
		}
	}

	// 如果没有获取到动态属性，使用静态属性作为备选
	if len(otherAttrDetails) == 0 {
		for _, attr := range otherAttributes {
			// 简单示例：为每个域名随机选择一些其他属性
			if strings.Contains(strings.ToLower(domainName), strings.ToLower(attr.AttributeValue)) {
				totalPriceFactor *= attr.PriceFactor
				totalGradeFactor += attr.GradeFactor
				otherAttrDetails = append(otherAttrDetails, model.AttributeDetail{
					Name:        attr.AttributeName,
					Value:       attr.AttributeValue,
					Description: attr.AttributeName,
					PriceFactor: attr.PriceFactor,
					GradeFactor: attr.GradeFactor,
				})
			}
		}
	}

	// 计算最终价格和等级
	finalPrice := basePrice * totalPriceFactor
	finalGrade := baseGrade + totalGradeFactor

	// 创建估价结果
	result := &model.EstimationResult{
		Domain:          domainName,
		Grade:           finalGrade,
		Price:           finalPrice,
		BaseAttributes:  baseAttrDetails,
		OtherAttributes: otherAttrDetails,
		EstimationDate:  time.Now(),
	}

	return result, nil
}

// parseDomain 解析域名，提取TLD、长度和结构等信息
func (s *DomainService) parseDomain(domainName string) (*model.Domain, error) {
	// 移除http://和https://前缀
	domainName = strings.TrimPrefix(domainName, "http://")
	domainName = strings.TrimPrefix(domainName, "https://")

	// 移除www.前缀
	domainName = strings.TrimPrefix(domainName, "www.")

	// 移除路径部分
	if idx := strings.Index(domainName, "/"); idx != -1 {
		domainName = domainName[:idx]
	}

	// 分割域名和TLD
	parts := strings.Split(domainName, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("无效的域名格式: %s", domainName)
	}

	// 获取TLD（最后一部分）
	tld := parts[len(parts)-1]

	// 获取域名主体（不含TLD）
	name := strings.Join(parts[:len(parts)-1], ".")

	// 确定域名结构
	structure := determineDomainStructure(name)

	// 尝试获取动态属性中的注册和到期日期
	registerDate := time.Now().AddDate(-1, 0, 0) // 默认：一年前注册
	expireDate := time.Now().AddDate(1, 0, 0)    // 默认：一年后到期

	// 尝试从动态属性服务获取WHOIS信息
	dynamicAttrs, err := s.dynamicAttrService.GetDynamicAttributes(domainName)
	if err == nil && dynamicAttrs != nil {
		// 解析注册日期
		if regDateStr, ok := dynamicAttrs["register_date"].(string); ok {
			if parsedDate, err := time.Parse("2006-01-02", regDateStr); err == nil {
				registerDate = parsedDate
			}
		}

		// 解析到期日期
		if expDateStr, ok := dynamicAttrs["expire_date"].(string); ok {
			if parsedDate, err := time.Parse("2006-01-02", expDateStr); err == nil {
				expireDate = parsedDate
			}
		}
	}

	// 创建域名对象
	domain := &model.Domain{
		Name:         domainName,
		TLD:          tld,
		Length:       len(name),
		Structure:    structure,
		RegisterDate: registerDate,
		ExpireDate:   expireDate,
	}

	return domain, nil
}

// determineDomainStructure 确定域名的结构类型
func determineDomainStructure(name string) string {
	// 检查是否为纯数字
	if regexp.MustCompile(`^\d+$`).MatchString(name) {
		return "纯数字"
	}

	// 检查是否为纯字母
	if regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(name) {
		return "纯字母"
	}

	// 检查是否为数字字母混合
	if regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(name) {
		return "数字字母混合"
	}

	// 检查是否包含连字符
	if strings.Contains(name, "-") {
		return "含连字符"
	}

	// 其他情况
	return "其他"
}
