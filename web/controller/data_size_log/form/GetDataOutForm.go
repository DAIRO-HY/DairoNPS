package form

type GetDataOutForm struct {

	/**
	 * 统计表标题列表
	 */
	Lables []string

	/**
	 * 入网流量
	 */
	InDatas []float64

	/**
	 * 出网流量
	 */
	OutDatas []float64

	/**
	 * 单位
	 */
	Unit string
}
