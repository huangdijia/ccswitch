<?php

declare(strict_types=1);
/**
 * This file is part of huangdijia/ccswitch.
 *
 * @link     https://github.com/huangdijia/ccswitch
 * @document https://github.com/huangdijia/ccswitch/blob/main/README.md
 * @contact  Your name <your-mail@gmail.com>
 */

namespace CCSwitch;

use Exception;

class Profiles
{
    public array $data = [];

    public function __construct(private string $path)
    {
        if (! file_exists($this->path)) {
            throw new Exception('Profiles file not found: ' . $this->path);
        }

        $content = file_get_contents($this->path);
        $this->data = json_decode($content, true);
    }

    public function getSettingsPath(): ?string
    {
        return $this->data['settingsPath'] ?? null;
    }

    public function has(string $name): bool
    {
        return array_key_exists($name, $this->data['profiles'] ?? []);
    }

    public function default(): string
    {
        return $this->data['default'] ?? 'default';
    }

    /**
     * @return array
     */
    public function get(string $name)
    {
        $env = $this->data['profiles'][$name] ?? [];
        $missing = ['ANTHROPIC_DEFAULT_HAIKU_MODEL', 'ANTHROPIC_DEFAULT_OPUS_MODEL', 'ANTHROPIC_DEFAULT_SONNET_MODEL', 'ANTHROPIC_SMALL_FAST_MODEL'];

        if (! isset($env['ANTHROPIC_MODEL'])) {
            return $env;
        }

        foreach ($missing as $missEnv) {
            if (! array_key_exists($missEnv, $env)) {
                $env[$missEnv] = $env['ANTHROPIC_MODEL'];
            }
        }

        return $env;
    }

    public function getLongOptions(): ?array
    {
        return array_map(
            fn ($key) => "{$key}::",
            array_keys($this->data['profiles'] ?? [])
        );
    }

    public function getUsage(?string $caller = null): string
    {
        $caller = $caller ?? 'php ' . basename(__FILE__);
        $opts = array_map(
            fn ($key) => "--{$key}",
            array_keys($this->data['profiles'] ?? [])
        );

        $usage = [];
        $usage[] = "Usage: {$caller} [" . implode(' | ', $opts) . '] [--settings=path/to/settings.json] [--profiles=path/to/ccs.json]';
        $usage[] = 'Options:';

        foreach ($this->data['descriptions'] ?? [] as $key => $description) {
            $usage[] = sprintf('  --%-6s %s', $key, $description);
        }

        return implode(PHP_EOL, $usage);
    }
}
