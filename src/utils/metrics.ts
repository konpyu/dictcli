export function calculateWER(reference: string, hypothesis: string): number {
  const refWords = reference.toLowerCase().split(/\s+/).filter(Boolean)
  const hypWords = hypothesis.toLowerCase().split(/\s+/).filter(Boolean)

  const m = refWords.length
  const n = hypWords.length

  if (m === 0) return n === 0 ? 0 : 1

  // Create DP table
  const dp: number[][] = Array(m + 1)
    .fill(0)
    .map(() => Array(n + 1).fill(0))

  // Initialize base cases
  for (let i = 0; i <= m; i++) dp[i][0] = i
  for (let j = 0; j <= n; j++) dp[0][j] = j

  // Fill DP table
  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      if (refWords[i - 1] === hypWords[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1]
      } else {
        dp[i][j] = Math.min(
          dp[i - 1][j] + 1, // deletion
          dp[i][j - 1] + 1, // insertion
          dp[i - 1][j - 1] + 1, // substitution
        )
      }
    }
  }

  return dp[m][n] / m
}

export function highlightDifferences(
  reference: string,
  hypothesis: string,
): Array<{ word: string; isError: boolean }> {
  const refWords = reference.split(/\s+/)
  const hypWords = hypothesis.split(/\s+/)

  const result: Array<{ word: string; isError: boolean }> = []

  const maxLen = Math.max(refWords.length, hypWords.length)
  for (let i = 0; i < maxLen; i++) {
    const refWord = refWords[i] || ''
    const hypWord = hypWords[i] || ''

    if (refWord.toLowerCase() === hypWord.toLowerCase()) {
      result.push({ word: refWord, isError: false })
    } else {
      if (refWord) result.push({ word: refWord, isError: true })
    }
  }

  return result
}
