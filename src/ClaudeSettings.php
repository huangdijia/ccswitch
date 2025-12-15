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

use stdClass;

/**
 * @property null|stdClass $env
 * @property null|string $model
 */
class ClaudeSettings
{
    private stdClass $data;

    public function __construct(private string $path)
    {
        $home = $_SERVER['HOME'] ?? ($_SERVER['HOMEDRIVE'] . $_SERVER['HOMEPATH']);
        $this->path = str_replace('~', $home, $path);

        if (! file_exists($this->path)) {
            file_put_contents($this->path, json_encode([], JSON_PRETTY_PRINT));
        }

        $this->data = $this->read();
    }

    public function __get(string $name)
    {
        return $this->data->{$name} ?? null;
    }

    public function __set(string $name, mixed $value)
    {
        $this->data->{$name} = $value;
    }

    public function __isset(string $name)
    {
        return isset($this->data->{$name});
    }

    public function __unset(string $name)
    {
        unset($this->data->{$name});
    }

    public function write(): void
    {
        $content = json_encode($this->data, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES | JSON_THROW_ON_ERROR);
        file_put_contents($this->path, $content);
    }

    private function read(): stdClass
    {
        $content = file_get_contents($this->path);
        return json_decode($content);
    }
}
