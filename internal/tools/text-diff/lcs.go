package textdiff

// op 是行级操作：kind 为 '='（相等）/ '-'（删除，仅旧侧）/ '+'（新增，仅新侧）。
// a/b 为各自侧的行索引（不适用时为 -1）。
type op struct {
	kind byte
	a, b int
}

// lcsDiff 通过 O(n·m) 动态规划 LCS 生成最短编辑脚本（与 Myers diff 结果等价：
// 均为最小编辑距离的删除+新增序列）。
// 回溯时优先保留相等行（稳定顺序）。输入行数受 maxLines 约束，DP 内存可控。
func lcsDiff(a, b []string) []op {
	n, m := len(a), len(b)
	if n == 0 {
		ops := make([]op, m)
		for j := 0; j < m; j++ {
			ops[j] = op{kind: '+', b: j}
		}
		return ops
	}
	if m == 0 {
		ops := make([]op, n)
		for i := 0; i < n; i++ {
			ops[i] = op{kind: '-', a: i}
		}
		return ops
	}

	// dp[i][j] = a[i:] 与 b[j:] 的 LCS 长度
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				dp[i][j] = dp[i+1][j+1] + 1
			} else if dp[i+1][j] >= dp[i][j+1] {
				dp[i][j] = dp[i+1][j]
			} else {
				dp[i][j] = dp[i][j+1]
			}
		}
	}

	ops := make([]op, 0, len(a)+len(b))
	i, j := 0, 0
	for i < n && j < m {
		if a[i] == b[j] {
			ops = append(ops, op{kind: '=', a: i, b: j})
			i++
			j++
		} else if dp[i+1][j] >= dp[i][j+1] {
			ops = append(ops, op{kind: '-', a: i})
			i++
		} else {
			ops = append(ops, op{kind: '+', b: j})
			j++
		}
	}
	for ; i < n; i++ {
		ops = append(ops, op{kind: '-', a: i})
	}
	for ; j < m; j++ {
		ops = append(ops, op{kind: '+', b: j})
	}
	return ops
}